import matplotlib.pyplot as plt

# Function to read data from a text file and parse it into a dictionary
def read_data_from_file(file_path):
    data = {
        "concurrent_users": [],
        "cpu_usage": [],
        "memory_usage": [],
        "execution_time": [],
        "success_rate": []
    }
    
    with open(file_path, 'r') as file:
        lines = file.readlines()
        for line in lines:
            if "Concurrent Users:" in line:
                data["concurrent_users"].append(int(line.split(":")[1].strip()))
            elif "Average CPU Usage:" in line:
                data["cpu_usage"].append(float(line.split(":")[1].strip().replace('%', '')))
            elif "Average Memory Usage:" in line:
                data["memory_usage"].append(float(line.split(":")[1].strip().replace(' MB', '')))
            elif "Average Total Execution Time:" in line:
                data["execution_time"].append(float(line.split(":")[1].strip().replace(' seconds', '')))
            elif "Average Success Rate:" in line:
                data["success_rate"].append(float(line.split(":")[1].strip().replace('%', '')))
    
    return data

# Reading data from files
microservice_data = read_data_from_file('output.txt')
monolith_data = read_data_from_file('micro-output.txt')

# Plotting CPU Usage
plt.figure(figsize=(7, 5))
plt.plot(monolith_data['concurrent_users'], monolith_data['cpu_usage'], label='Monolith', marker='o')
plt.plot(microservice_data['concurrent_users'], microservice_data['cpu_usage'], label='Microservice', marker='o')
plt.title('CPU Usage vs Concurrent Users')
plt.xlabel('Concurrent Users')
plt.ylabel('CPU Usage (%)')
plt.legend()
plt.tight_layout()
plt.savefig('cpu_usage_vs_concurrent_users.png')
plt.close()

# Plotting Memory Usage
plt.figure(figsize=(7, 5))
plt.plot(monolith_data['concurrent_users'], monolith_data['memory_usage'], label='Monolith', marker='o')
plt.plot(microservice_data['concurrent_users'], microservice_data['memory_usage'], label='Microservice', marker='o')
plt.title('Memory Usage vs Concurrent Users')
plt.xlabel('Concurrent Users')
plt.ylabel('Memory Usage (MB)')
plt.legend()
plt.tight_layout()
plt.savefig('memory_usage_vs_concurrent_users.png')
plt.close()

# Plotting Total Execution Time
plt.figure(figsize=(7, 5))
plt.plot(monolith_data['concurrent_users'], monolith_data['execution_time'], label='Monolith', marker='o')
plt.plot(microservice_data['concurrent_users'], microservice_data['execution_time'], label='Microservice', marker='o')
plt.title('Total Execution Time vs Concurrent Users')
plt.xlabel('Concurrent Users')
plt.ylabel('Total Execution Time (seconds)')
plt.legend()
plt.tight_layout()
plt.savefig('total_execution_time_vs_concurrent_users.png')
plt.close()

# Plotting Success Rate
plt.figure(figsize=(7, 5))
plt.plot(monolith_data['concurrent_users'], monolith_data['success_rate'], label='Monolith', marker='o')
plt.plot(microservice_data['concurrent_users'], microservice_data['success_rate'], label='Microservice', marker='o')
plt.title('Success Rate vs Concurrent Users')
plt.xlabel('Concurrent Users')
plt.ylabel('Success Rate (%)')
plt.legend()
plt.tight_layout()
plt.savefig('success_rate_vs_concurrent_users.png')
plt.close()
