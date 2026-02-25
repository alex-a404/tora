import numpy as np

DEBUG = True

def log(*args):
    if DEBUG:
        print(*args)

def get_route_cost(route, dist_matrix, node_depths, direction, penalty_weight=500):
    log("\n[get_route_cost] Route:", route)
    if not route:
        return 0

    # FIX 2: Use the Node IDs directly as indices!
    route_idx = np.array(route)

    # --- Distance cost ---
    total_dist = np.sum(dist_matrix[route_idx[:-1], route_idx[1:]])
    log("  Distance cost:", total_dist)

    # --- Corridor penalty ---
    depths = np.array([node_depths[node] for node in route])
    diffs = np.diff(depths)

    backtrack_mask = np.sign(diffs) != direction
    backtrack_penalty = np.sum(backtrack_mask * np.abs(diffs) * penalty_weight)

    total_cost = total_dist + backtrack_penalty
    log("  TOTAL COST:", total_cost)

    return total_cost

def is_precedence_valid(route, dependencies):
    log("\n[is_precedence_valid] Checking route:", route)
    pos = {node_id: i for i, node_id in enumerate(route)}
    log("  Position map:", pos)

    for p, d in dependencies:
        if p not in pos or d not in pos:
            log(f"  WARNING: dependency ({p},{d}) missing in route")
            continue
        if pos[p] >= pos[d]:
            log(f"  ❌ Invalid precedence: pickup {p} at {pos[p]} after dropoff {d} at {pos[d]}")
            return False

    log("  ✅ Precedence valid")
    return True

def optimize_itinerary(dist_matrix, node_depths, current_route,
                       new_request, dependencies, direction, locked_count=1):

    log("\n==============================")
    log("[optimize_itinerary] START")
    log("Current route:", current_route)

    p_new, d_new = new_request
    best_route = None
    min_cost = float('inf')

    # --- STAGE 1: CHEAPEST INSERTION ---
    for i in range(locked_count, len(current_route) + 1):
        for j in range(i + 1, len(current_route) + 2):
            temp_route = list(current_route)
            temp_route.insert(i, p_new)
            temp_route.insert(j, d_new)

            if is_precedence_valid(temp_route, dependencies):
                # Removed node_id_to_idx
                cost = get_route_cost(temp_route, dist_matrix, node_depths, direction)
                if cost < min_cost:
                    min_cost = cost
                    best_route = temp_route

    # --- STAGE 2: 2-OPT ---
    if best_route and len(best_route) > 3:
        best_route = run_2opt(best_route, dist_matrix, node_depths,
                              dependencies, direction, locked_count)

    return best_route

def run_2opt(route, dist_matrix, node_depths, dependencies, direction, locked_count):
    best_route = list(route)
    best_cost = get_route_cost(best_route, dist_matrix, node_depths, direction)

    improved = True
    while improved:
        improved = False
        for i in range(locked_count, len(best_route) - 1):
            for j in range(i + 1, len(best_route)):
                new_route = best_route[:i] + best_route[i:j+1][::-1] + best_route[j+1:]

                if is_precedence_valid(new_route, dependencies):
                    new_cost = get_route_cost(new_route, dist_matrix, node_depths, direction)
                    if new_cost < best_cost:
                        best_cost = new_cost
                        best_route = new_route
                        improved = True
    return best_route