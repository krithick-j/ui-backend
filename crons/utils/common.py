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
        print(f"Error in get_personal_rs: {str(e)}")
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
        SELECT SUM(direct_bv) FROM rsp_transactions WHERE distrib_id ='{distrib_id}'
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

        return count
    
    except Exception as e:
        print(f"Error in get_step_by_distrib_id: {str(e)}")
        return 0


    
