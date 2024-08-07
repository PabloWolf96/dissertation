#!/bin/bash

# Base URL of your application
BASE_URL="http://localhost:8000"

# Function to generate a random number between min and max
random() {
    min=$1
    max=$2
    echo $((RANDOM % (max - min + 1) + min))
}

# Function to measure time and extract metrics
measure_time() {
    start_time=$(date +%s.%N)
    output=$("$@")
    end_time=$(date +%s.%N)
    duration=$(echo "$end_time - $start_time" | bc)
    echo "$duration:$output"
}

# Initialize variables
login_time=0
avg_cart_time=0
max_cart_time=0
checkout_time=0
total_time=0
num_products=0
http_code=0
cart_http_code=0
checkout_http_code=0

# Login and get token
login_response=$(measure_time curl -s -w "%{http_code}" -X POST "$BASE_URL/login" \
     -H "Content-Type: application/json" \
     -d '{"username":"owen","password":"12345"}')

login_time=$(echo $login_response | cut -d':' -f1)
login_output=$(echo $login_response | cut -d':' -f2-)
http_code=$(echo $login_output | sed 's/.*\(...\)$/\1/')
token=$(echo $login_output | sed 's/...$//' | jq -r '.token')

if [ -z "$token" ] || [ "$http_code" != "200" ]; then
    echo "Failed to login. HTTP Code: $http_code, Response: $login_output"
else
    echo "Login time: $login_time seconds"

    # Add random number of products to cart
    num_products=$(random 1 10)
    total_cart_time=0
    max_cart_time=0

    for i in $(seq 1 $num_products); do
        product_id=$(random 1 1000)
        quantity=$(random 1 5)
        
        add_to_cart_response=$(measure_time curl -s -w "%{http_code}" -X POST "$BASE_URL/cart" \
             -H "Content-Type: application/json" \
             -H "Authorization: Bearer $token" \
             -d "{\"product_id\":$product_id,\"quantity\":$quantity}")
        
        cart_time=$(echo $add_to_cart_response | cut -d':' -f1)
        cart_output=$(echo $add_to_cart_response | cut -d':' -f2-)
        cart_http_code=$(echo $cart_output | sed 's/.*\(...\)$/\1/')
        
        total_cart_time=$(echo "$total_cart_time + $cart_time" | bc)
        max_cart_time=$(echo "if ($cart_time > $max_cart_time) $cart_time else $max_cart_time" | bc)
        
        echo "Added product $product_id, quantity $quantity to cart. Time: $cart_time seconds, HTTP Code: $cart_http_code"
    done

    avg_cart_time=$(echo "scale=4; $total_cart_time / $num_products" | bc)
    echo "Average add to cart time: $avg_cart_time seconds"
    echo "Max add to cart time: $max_cart_time seconds"

    # Perform checkout
    checkout_response=$(measure_time curl -s -w "%{http_code}" -X POST "$BASE_URL/cart/checkout" \
         -H "Content-Type: application/json" \
         -H "Authorization: Bearer $token")

    checkout_time=$(echo $checkout_response | cut -d':' -f1)
    checkout_output=$(echo $checkout_response | cut -d':' -f2-)
    checkout_http_code=$(echo $checkout_output | sed 's/.*\(...\)$/\1/')

    echo "Checkout time: $checkout_time seconds, HTTP Code: $checkout_http_code"

    # Calculate total execution time
    total_time=$(echo "$login_time + $total_cart_time + $checkout_time" | bc)
    echo "Total execution time: $total_time seconds"
fi

# Output in a format easy to parse later, even if there were errors
echo "METRICS:$login_time:$avg_cart_time:$max_cart_time:$checkout_time:$total_time:$num_products:$http_code:$cart_http_code:$checkout_http_code"