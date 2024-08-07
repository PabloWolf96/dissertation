import re
from collections import defaultdict

def safe_float(value):
    try:
        return float(value)
    except (ValueError, TypeError):
        return 0.0

def parse_load_test_file(file_path):
    with open(file_path, 'r') as file:
        content = file.read()

    pattern = re.compile(r'Final Results for (\d+) concurrent users:\n(.*?)\n---------------------------------------', re.DOTALL)
    matches = pattern.findall(content)

    results = defaultdict(list)

    for match in matches:
        users = int(match[0])
        data = match[1].strip().split('\n')
        
        result = {
            'users': users,
            'login_time': safe_float(re.search(r'Average Login Time: ([-\d.]+)', '\n'.join(data)).group(1) if re.search(r'Average Login Time: ([-\d.]+)', '\n'.join(data)) else 0),
            'add_to_cart_time': safe_float(re.search(r'Average Add to Cart Time: ([-\d.]+)', '\n'.join(data)).group(1) if re.search(r'Average Add to Cart Time: ([-\d.]+)', '\n'.join(data)) else 0),
            'max_add_to_cart_time': safe_float(re.search(r'Max Add to Cart Time: ([-\d.]+)', '\n'.join(data)).group(1) if re.search(r'Max Add to Cart Time: ([-\d.]+)', '\n'.join(data)) else 0),
            'checkout_time': safe_float(re.search(r'Average Checkout Time: ([-\d.]+)', '\n'.join(data)).group(1) if re.search(r'Average Checkout Time: ([-\d.]+)', '\n'.join(data)) else 0),
            'total_execution_time': safe_float(re.search(r'Average Total Execution Time: ([-\d.]+)', '\n'.join(data)).group(1) if re.search(r'Average Total Execution Time: ([-\d.]+)', '\n'.join(data)) else 0),
            'products_per_cart': safe_float(re.search(r'Average Products per Cart: ([-\d.]+)', '\n'.join(data)).group(1) if re.search(r'Average Products per Cart: ([-\d.]+)', '\n'.join(data)) else 0),
            'success_rate': safe_float(re.search(r'Success Rate: ([-\d.]+)', '\n'.join(data)).group(1) if re.search(r'Success Rate: ([-\d.]+)', '\n'.join(data)) else 0),
            'cpu_usage': safe_float(re.search(r'CPU Usage: ([-\d.]+)', '\n'.join(data)).group(1) if re.search(r'CPU Usage: ([-\d.]+)', '\n'.join(data)) else 0),
            'memory_usage': int(re.search(r'Memory Usage: (\d+)', '\n'.join(data)).group(1) if re.search(r'Memory Usage: (\d+)', '\n'.join(data)) else 0)
        }
        
        results[users].append(result)

    return results

def calculate_averages(results):
    averages = {}

    for users, data in results.items():
        avg = {
            'users': users,
            'login_time': sum(d['login_time'] for d in data) / len(data),
            'add_to_cart_time': sum(d['add_to_cart_time'] for d in data) / len(data),
            'max_add_to_cart_time': sum(d['max_add_to_cart_time'] for d in data) / len(data),
            'checkout_time': sum(d['checkout_time'] for d in data) / len(data),
            'total_execution_time': sum(d['total_execution_time'] for d in data) / len(data),
            'products_per_cart': sum(d['products_per_cart'] for d in data) / len(data),
            'success_rate': sum(d['success_rate'] for d in data) / len(data),
            'cpu_usage': sum(d['cpu_usage'] for d in data) / len(data),
            'memory_usage': sum(d['memory_usage'] for d in data) / len(data)
        }
        averages[users] = avg

    return averages

def write_averages_to_file(averages, output_file):
    with open(output_file, 'w') as file:
        file.write("Averages for each batch of concurrent users:\n")
        file.write("-" * 50 + "\n")
        for users, avg in sorted(averages.items()):
            file.write(f"Concurrent Users: {users}\n")
            file.write(f"  Average Login Time: {avg['login_time']:.4f} seconds\n")
            file.write(f"  Average Add to Cart Time: {avg['add_to_cart_time']:.4f} seconds\n")
            file.write(f"  Average Max Add to Cart Time: {avg['max_add_to_cart_time']:.4f} seconds\n")
            file.write(f"  Average Checkout Time: {avg['checkout_time']:.4f} seconds\n")
            file.write(f"  Average Total Execution Time: {avg['total_execution_time']:.4f} seconds\n")
            file.write(f"  Average Products per Cart: {avg['products_per_cart']:.2f}\n")
            file.write(f"  Average Success Rate: {avg['success_rate']:.2f}%\n")
            file.write(f"  Average CPU Usage: {avg['cpu_usage']:.2f}%\n")
            file.write(f"  Average Memory Usage: {avg['memory_usage']:.2f} MB\n")
            file.write("-" * 50 + "\n")
    print(f"Results have been written to {output_file}")

def main():
    input_file = 'results.txt'  # Replace with your input file path
    output_file = 'output.txt'  # Output file path
    results = parse_load_test_file(input_file)
    averages = calculate_averages(results)
    write_averages_to_file(averages, output_file)

if __name__ == "__main__":
    main()