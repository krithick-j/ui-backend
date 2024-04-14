from datetime import datetime
from typing import List

from utils.config import DB, BV_VALUE

def get_all_users_with_frequency() -> List[str]:
    '''
    '''
    users = []
    with DB.cursor() as cursor:
        sql = """
        SELECT u.distrib_id, cf.frequency FROM users u
        LEFT JOIN cheque_frequencies cf ON u.distrib_id = cf.distrib_id
        """
        cursor.execute(sql)
        rows = cursor.fetchall()
        for row in rows:
            users.append({"id":row[0], "frequency":row[1]})
    
    return users

def get_rank(user_id:str) -> str:
    ''' Calculate rank here
    The logic is for current month how many 
        1. Cheques collected
        2. How many sign up happened
        3. How many bv they generated
    '''
    return '2.0'

def alluser_upsert_rank() -> None:
    cursor = DB.cursor()
    for user in users:
        rank = get_rank(user["id"])
        sql = f"""
            INSERT INTO user_rank (user_id, year, month, rank)
            VALUES ('{user["id"]}', {current_dt.year}, {current_dt.month}, '{rank}')
            ON DUPLICATE KEY UPDATE rank = '{rank}';
        """
        cursor.execute(sql)
    DB.commit()
    cursor.close()

def process_ep(distrib_id: str, 
               place: str) -> None:
    print(f"Processing ev for {distrib_id} on {place}")
    cursor = DB.cursor()
    sqlInsert = f"""
            INSERT INTO ep_transactions 
              (distrib_id, reference, value) VALUES
              ('{distrib_id}', '{place}', {BV_VALUE*2});
        """
    sqlDeduct = f"""
            INSERT INTO bv_transactions (distrib_id, place, side, trans_type, bv_value) VALUES
            ('{distrib_id}', '{place}', 'left', 'redeem', -{BV_VALUE}),
            ('{distrib_id}', '{place}', 'right', 'redeem', -{BV_VALUE});
        """
    sqlUpdate = f"""
            UPDATE cheque_frequencies 
                SET frequency = frequency + 1 
            WHERE distrib_id = '{distrib_id}'
        """
    # Insert into ep_transaction
    cursor.execute(sqlInsert, multi=True)
    # Deduct bv from tracking center
    cursor.execute(sqlDeduct, multi=True)
    # Reset the frequencies
    cursor.execute(sqlUpdate, multi=True)

    DB.commit()
    cursor.close()

def allusers_update_earning_point() -> None:
    cursor = DB.cursor()
    for user in users:
        if user["frequency"] is None or (user["frequency"] % 5 != 0):
            # Skipping earning point check for {user['id']}
            continue
        sql = f"""
            SELECT 
              place,
              SUM(IF(side='left', bv_value, 0)) as left_point,
              SUM(IF(side='right', bv_value, 0)) as right_point
            FROM bv_transactions 
            WHERE distrib_id = '{user["id"]}'
            GROUP BY place;
        """
        cursor.execute(sql)
        rows  = cursor.fetchall()
        for row in rows:
            if row[1] >= BV_VALUE and row[2] >= BV_VALUE:
                process_ep(user['id'], row[0])
                break
        # if user.left < 4000 or user.rightbv < 4000:
            # continue
        # Deduct 400 from left and deduct 4000 from right
        # add earning_point to the user
        # reset the frequency to 0

def activate_bv() -> None:
    cursor = DB.cursor()
    sql = """
    update bv_transactions set is_active = 1 where DATE(activate_date) = DATE(NOW())
    """
    cursor.execute(sql)
    DB.commit()
    cursor.close()
    return None

def main()->None:
    # Update the rank for all users
    # alluser_upsert_rank()
    # Check earning point on 6th check
    allusers_update_earning_point()
    # Close the db at the end
    # activate_bv()
    DB.close()

#Keep the users in global scope to get repetive db hit
current_dt = datetime.now()
users = get_all_users_with_frequency()
main()
