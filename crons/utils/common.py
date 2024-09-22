from typing import Dict, List, Tuple
from utils.config import DB
from datetime import datetime, timedelta


def get_all_users_with_frequency() -> List[str]:
    '''
    '''
    users = []
    with DB.cursor() as cursor:
        sql = """
        SELECT u.distrib_id, cf.frequency, u.highest_rank, u.current_rank FROM users u
        LEFT JOIN cheque_frequencies cf ON u.distrib_id = cf.distrib_id
        """
        cursor.execute(sql)
        rows = cursor.fetchall()
        for row in rows:
            users.append({"id":row[0], "frequency":row[1], "highest_rank": row[2], "current_rank": row[3]})
        cursor.close()
        
    return users

def get_rank_goal()->Dict:
    rank = {}
    cursor = DB.cursor()
    sql = """
    SELECT rank_id, type, target from rank_details;
    """
    cursor.execute(sql)
    rows = cursor.fetchall()
    for row in rows:
        if not rank.get(row[0]):
            rank[row[0]] = {}
        rank[row[0]][row[1]] =  row[2]
    cursor.close()
    return rank

def get_personal_rsp(distrib_id: str) -> float:
    try:
        cursor = DB.cursor()
        sql = f"""
        SELECT COALESCE(SUM(rsp), 0) FROM rsp_transactions WHERE distrib_id ='{distrib_id}'
        """
        cursor.execute(sql)
        result = cursor.fetchone()
        cursor.close()
        return result[0] if result else 0.0

    except Exception as e:
        print(f"Error in get_personal_rsp: {str(e)}")
        return 0.0

def get_group_rsp_by_distrib_id(distrib_id: str) -> Tuple[float, int]:
    rsp_sum = 0.0
    try:
        # Get the referred distrib IDs
        referred_distrib_ids = get_referred_users_by_distrib_id(distrib_id)
        if not referred_distrib_ids:
            return rsp_sum, 200

        for referred_distrib_id in referred_distrib_ids:
            # Prevent infinite loop if distrib_id and referred_distrib_id are the same
            if distrib_id == referred_distrib_id:
                continue

            # Get the RSP sum for the referred distrib ID
            total_rsp_for_one_distrib = get_personal_rsp(referred_distrib_id)
            rsp_sum += total_rsp_for_one_distrib

            # Recursively get RSP sum for the group
            recursive_rsp_sum, status = get_group_rsp_by_distrib_id(referred_distrib_id)
            if status != 200:
                return recursive_rsp_sum, status

            rsp_sum += recursive_rsp_sum

        return rsp_sum, 200

    except Exception as e:
        print(f"Error in get_group_rsp_by_distrib_id: {str(e)}")
        return rsp_sum, 500

def get_referred_users_by_distrib_id(distrib_id: str) -> list:
    try:
        cursor = DB.cursor()
        sql = f"""
        SELECT distrib_id FROM users WHERE ref_distrib_id ='{distrib_id}'
        """
        cursor.execute(sql)
        result = cursor.fetchall()
        cursor.close()
        return [row[0] for row in result]

    except Exception as e:
        print(f"Error in get_referred_users_by_distrib_id: {str(e)}")
        return []
    
def get_direct_bv(distrib_id: str) -> float:
    try:
        cursor = DB.cursor()
        sql = f"""
        SELECT COALESCE(SUM(direct_bv), 0) FROM rsp_transactions WHERE distrib_id ='{distrib_id}'
        """
        cursor.execute(sql)
        result = cursor.fetchone()
        cursor.close()
        return result[0] if result else 0.0

    except Exception as e:
        print(f"Error in get_direct_bv: {str(e)}")
        return 0.0

def get_step_by_distrib_id(distrib_id: str) -> int:
    try:
        cursor = DB.cursor()
        # Calculate the first and last days of the current month
        first_of_month = datetime(datetime.now().year, datetime.now().month, 1)
        last_of_month = first_of_month + timedelta(days=32)
        last_of_month = datetime(last_of_month.year, last_of_month.month, 1) - timedelta(seconds=1)

        # SQL query to count the transactions
        query = """
            SELECT COUNT(*)
            FROM bv_transactions
            WHERE distrib_id = %s
            AND trans_type = %s
            AND created_at BETWEEN %s AND %s
        """
        
        # Execute the query with the provided parameters
        cursor.execute(query, (distrib_id, 'cheque', first_of_month, last_of_month))
        count = cursor.fetchone()[0]
        count = count//2
        return count
    
    except Exception as e:
        print(f"Error in get_step_by_distrib_id: {str(e)}")
        return 0


def get_current_rank_value_by_distrib_id(distrib_id: str):
    cursor = DB.cursor()
    query = '''
    SELECT current_rank
    FROM users
    WHERE distrib_id = %s
    '''
    cursor.execute(query, (distrib_id,))
    result = cursor.fetchone()
    if result:
        return result[0]
    else:
        raise ValueError(f"No rank found for distrib_id {distrib_id}")
    
def get_current_rank_arr_by_distrib_id(referred_distrib_ids: list[str]):
    cursor = DB.cursor()
    query = '''
    SELECT current_rank
    FROM users
    WHERE distrib_id IN (%s)
    ''' % ','.join(['%s'] * len(referred_distrib_ids))
    cursor.execute(query, referred_distrib_ids)
    referred_ranks = [float(row[0]) for row in cursor.fetchall()]
    return referred_ranks

def get_referred_distrib_id(distrib_id: str):
    query = '''
    SELECT DISTINCT distrib_id
    FROM users
    WHERE ref_distrib_id = %s AND distrib_id != ref_distrib_id
    '''
    cursor = DB.cursor()
    cursor.execute(query, (distrib_id,))
    referred_distributors = [row[0] for row in cursor.fetchall()]
    return referred_distributors

def get_referral_chain_by_distrib_id(distrib_id):
    user_arr = []

    def build_referral_chain(distrib_id):
        referred_distrib_ids = get_referred_distrib_id(distrib_id)

        for distrib_id in referred_distrib_ids:
            user_arr.append(distrib_id)  # Add to the chain
            err = build_referral_chain(distrib_id)  # Recurse to find the next level of referrals
            if err is not None:
                return err  # Propagate the error if something went wrong during recursion

        return None

    # Start building the referral chain
    err = build_referral_chain(distrib_id)
    if err is not None:
        return str(err)

    # Return the referral chain as part of the response
    return user_arr

def get_group_performance_by_distrib_id(distrib_id: str):
    # Get the referral chain for the given distrib_id
    referred_distrib_ids = get_referral_chain_by_distrib_id(distrib_id)
    # Get the current rank of the original distributor
    current_rank_distrib_id = get_current_rank_value_by_distrib_id(distrib_id)

    # Count the number of referred distributors with rank >= the original distributor
    count = 0
    if referred_distrib_ids:
        referred_ranks = get_current_rank_arr_by_distrib_id(referred_distrib_ids)
        # Count the number of ranks that are >= current_rank_distrib_id
        count = sum(1 for rank in referred_ranks if rank >= current_rank_distrib_id)

    return count
    
