#!/bin/bash

BASE_LOAD=100
INCREMENTS=(10 20 30 40 50 60 70 80 90 100)
DURATION=30  # Increased to 5 minutes (300 seconds) per increment

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
    
    # Array to store results
    declare -a results

    # Run complex test scenario
    while [ $(($(date +%s) - increment_start_time)) -lt $DURATION ]; do
        for i in $(seq 1 $CURRENT_LOAD); do
            output=$(./complex_test.sh) &  # Run in background to simulate true concurrency
            results+=("$output")
        done
        sleep 1  # Short pause to prevent overwhelming the system

        # Collect and log intermediate metrics every minute
        if [ $(( ($(date +%s) - increment_start_time) % 60 )) -eq 0 ]; then
            system_metrics=$(collect_system_metrics)
            IFS=':' read -r cpu_usage memory_usage <<< "$system_metrics"
            log_results "Intermediate metrics at $(date +%M:%S) - CPU Usage: $cpu_usage%, Memory Usage: $memory_usage MB"
        fi
    done

    # Wait for all background processes to finish
    wait

    # Process and aggregate results
    # ... (rest of the result processing code remains the same)

    log_results "Final Results for $CURRENT_LOAD concurrent users:"
    log_results "Average Login Time: $avg_login_time seconds"
    log_results "Average Add to Cart Time: $avg_cart_time seconds"
    log_results "Max Add to Cart Time: $total_max_cart_time seconds"
    log_results "Average Checkout Time: $avg_checkout_time seconds"
    log_results "Average Total Execution Time: $avg_execution_time seconds"
    log_results "Average Products per Cart: $avg_products"
    log_results "Success Rate: $success_rate%"
    log_results "CPU Usage: $cpu_usage%"
    log_results "Memory Usage: $memory_usage MB"
    log_results "---------------------------------------"

    # Allow system to stabilize before next increment
    log_results "Cooling down for 2 minutes before next increment"
    sleep 30
done

log_results "Extended load test completed"