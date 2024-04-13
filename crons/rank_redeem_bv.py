from datetime import datetime
from typing import List

from utils.config import DB

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

def allusers_update_earning_point() -> None:
    cursor = DB.cursor()
    for user in users:
        if user["frequency"] is None or user["frequency"] < 5:
            print(f"Skipping earning point check for {user['id']}")
            continue
        print(f"further checking for {user['id']}")
        for center_code in ['001', '002', '003']:
            sql = f"""
                SELECT left_point, right_point From tracking_centers where distrib_id = {user['id']} AND place = {center_code};
            """
            cursor.execute(sql)
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
    # allusers_update_earning_point()
    # Close the db at the end
    activate_bv()
    DB.close()

#Keep the users in global scope to get repetive db hit
current_dt = datetime.now()
users = get_all_users_with_frequency()
main()
