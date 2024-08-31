from utils.common import (
    get_all_users_with_frequency,
    get_direct_bv,
    get_group_performance_by_distrib_id,
    get_group_rsp_by_distrib_id,
    get_personal_rsp,
    get_rank_goal,
    get_step_by_distrib_id,
)
from utils.config import DB

def check_personal_rsp(personal_rsp, rank: int)-> bool:
    return personal_rsp >= GOAL[rank]["PRSP"]

def check_group_rsp(group_rsp, rank: int)-> bool:
    return group_rsp >= GOAL[rank]["GRSP"]

def check_direct_bv(direct_bv, rank: int)-> bool:
    return direct_bv >= GOAL[rank]["DRBV"]

def check_step_value(step_value, rank: int)-> bool:
    return step_value >= GOAL[rank]["STEP"]

def check_group_performance(grp_performance, rank: int) -> bool:
    return grp_performance >= GOAL[rank]["GPRF"]

def check_rank_for_user(distrib_id: str)->int:
    orank = [2, 2.5, 3, 3.5, 4]
    ranks = [2.5, 3, 3.5, 4]
    grp_performance = get_group_performance_by_distrib_id(distrib_id)
    step_value = get_step_by_distrib_id(distrib_id)
    direct_bv = get_direct_bv(distrib_id)
    group_rsp, _ = get_group_rsp_by_distrib_id(distrib_id)
    personal_rsp = get_personal_rsp(distrib_id)
    
    checks = {
        rank: all([
            check_personal_rsp(personal_rsp, rank),
            check_group_rsp(group_rsp, rank),
            check_direct_bv(direct_bv, rank),
            check_step_value(step_value, rank),
            check_group_performance(grp_performance, rank)
        ])
        for rank in ranks
    }   
         
    for idx, rank in enumerate(ranks):
        if not checks[rank]:
            return orank[idx]

def check_rank_for_all()->None:
   for user in users:
       highest_rank = user['highest_rank']
       new_rank = check_rank_for_user(user['id'])
       update_current_highest_rank(user['id'], new_rank, highest_rank)
       detect_direct_bv(user['id'])
       
def detect_direct_bv(distribID)->None:
    direct_bv = get_direct_bv(distribID)
    cursor = DB.cursor()
    if direct_bv >= 500:
        query = f"INSERT INTO rsp_transactions (distrib_id, direct_bv) VALUES ('{distribID}', -500)"
        cursor.execute(query)
        DB.commit()
        cursor.close()

def update_current_highest_rank(distrib_id, new_rank, highest_rank):
    try:
       cursor = DB.cursor()
       params = [new_rank]

       sql = "UPDATE users SET current_rank = %s"

        # Adding an additional update condition if the new rank is higher
       if new_rank > highest_rank:
           sql += ", highest_rank = %s"
           params.append(new_rank)
           

        # Completing the query with the WHERE clause
       sql += " WHERE distrib_id = %s"
       params.append(distrib_id)

        # Executing the query with the parameters
       cursor.execute(sql, params)
       DB.commit()
       cursor.close()
        
    except Exception as e:
        print(f"Error in update_current_highest_rank: {str(e)}")
        return 0.0
    
        
users = get_all_users_with_frequency()
GOAL = get_rank_goal()
check_rank_for_all()
# distrib_id = 'IN-00001'
# GOAL = get_rank_goal()
# print("bronze",check_rank_for_user(distrib_id, 1, 2))
# print("silver",check_rank_for_user(distrib_id, 2, 2))
# print("gold",check_rank_for_user(distrib_id, 3, 2))
# print("sapphire",check_rank_for_user(distrib_id, 4, 2))
# print("platnstar",check_rank_for_user(distrib_id, 5, 2))
