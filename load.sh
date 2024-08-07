#!/bin/bash

BASE_LOAD=100
INCREMENTS=(10 10 10 20 20 20 30 30 30 40 40 40 50 50 50 60 60 60 70 70 70 80 80 80 90 90 90 100 100 100 110 110 110 120 120 120 130 130 130 140 140 140 150 150 150 160 160 160)
DURATION=300  # 5 minutes per increment

start_time=$(date +%s)

# Function to log results
log_results() {
    echo "$1" | tee -a results.txt
}

# Initialize results file
echo "Extended Load Test Results - $(date)" > results.txt
echo "===============================" >> results.txt

collect_system_metrics() {
    cpu_usage=$(top -bn1 | grep "Cpu(s)" | sed "s/.*, *\([0-9.]*\)%* id.*/\1/" | awk '{print 100 - $1}')
    memory_usage=$(free -m | awk '/Mem:/ {print $3}')
    echo "$cpu_usage:$memory_usage"
}

for INCREMENT in "${INCREMENTS[@]}"; do
    CURRENT_LOAD=$((BASE_LOAD * INCREMENT / 100))
    log_results "Increasing load to $CURRENT_LOAD concurrent users (${INCREMENT}% of base load)"

    increment_start_time=$(date +%s)
    
    # Run complex test scenario
    for i in $(seq 1 $CURRENT_LOAD); do
        ./complex_test.sh > "temp_result_$i.txt" &
    done

    # Wait for all background processes to finish
    wait

    # Process and aggregate results
    total_login_time=0
    total_avg_cart_time=0
    total_max_cart_time=0
    total_checkout_time=0
    total_execution_time=0
    total_products=0
    success_count=0
    valid_results=0

    for i in $(seq 1 $CURRENT_LOAD); do
        result=$(<"temp_result_$i.txt")
        if [[ $result == *"METRICS:"* ]]; then
            IFS=':' read -r -a metrics <<< "${result#*METRICS:}"
            if [[ ${#metrics[@]} -eq 9 ]]; then
                total_login_time=$(echo "$total_login_time + ${metrics[0]}" | bc)
                total_avg_cart_time=$(echo "$total_avg_cart_time + ${metrics[1]}" | bc)
                total_max_cart_time=$(echo "if (${metrics[2]} > $total_max_cart_time) ${metrics[2]} else $total_max_cart_time" | bc)
                total_checkout_time=$(echo "$total_checkout_time + ${metrics[3]}" | bc)
                total_execution_time=$(echo "$total_execution_time + ${metrics[4]}" | bc)
                total_products=$(echo "$total_products + ${metrics[5]}" | bc)
                if [ "${metrics[6]}" = "200" ] && [ "${metrics[7]}" = "200" ] && [ "${metrics[8]}" = "200" ]; then
                    success_count=$((success_count + 1))
                fi
                valid_results=$((valid_results + 1))
            else
                log_results "Invalid metric count in result $i: ${#metrics[@]}"
            fi
        else
            log_results "No METRICS found in result $i"
        fi
        rm "temp_result_$i.txt"
    done

    log_results "Final Results for $CURRENT_LOAD concurrent users:"
    if [ $valid_results -gt 0 ]; then
        avg_login_time=$(echo "scale=4; $total_login_time / $valid_results" | bc)
        avg_cart_time=$(echo "scale=4; $total_avg_cart_time / $valid_results" | bc)
        avg_checkout_time=$(echo "scale=4; $total_checkout_time / $valid_results" | bc)
        avg_execution_time=$(echo "scale=4; $total_execution_time / $valid_results" | bc)
        avg_products=$(echo "scale=2; $total_products / $valid_results" | bc)
        success_rate=$(echo "scale=2; $success_count * 100 / $valid_results" | bc)

        log_results "Average Login Time: $avg_login_time seconds"
        log_results "Average Add to Cart Time: $avg_cart_time seconds"
        log_results "Max Add to Cart Time: $total_max_cart_time seconds"
        log_results "Average Checkout Time: $avg_checkout_time seconds"
        log_results "Average Total Execution Time: $avg_execution_time seconds"
        log_results "Average Products per Cart: $avg_products"
        log_results "Success Rate: $success_rate%"
    else
        log_results "No valid results collected for this increment"
    fi

    # Collect final system metrics
    system_metrics=$(collect_system_metrics)
    IFS=':' read -r cpu_usage memory_usage <<< "$system_metrics"
    log_results "CPU Usage: $cpu_usage%"
    log_results "Memory Usage: $memory_usage MB"
    log_results "---------------------------------------"

    # Allow system to stabilize before next increment
    log_results "Cooling down for 30 seconds before next increment"
    sleep 30
done

log_results "Extended load test completed"