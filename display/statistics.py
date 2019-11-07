import ifaddr
import subprocess

# Gigabyte to kibibyte conversion factor
KIB_GB = 976562.5


# Get the disk statistics
def disk(mount_point):
    # Get disk usage status
    raw = subprocess.run(["df", "--output=used,pcent,avail,size", mount_point], stdout=subprocess.PIPE, stderr=subprocess.PIPE)

    # Check return code
    if raw.returncode != 0:
        raise ChildProcessError(f"non-zero exit code for df: {raw.stderr.decode()}")

    # Remove empty values
    parsed = raw.stdout.decode().split("\n")[1].split(" ")
    for i, value in enumerate(parsed):
        if value == "":
            del parsed[i]

    return {
        "used": int(parsed[0]) / KIB_GB,
        "percent": parsed[1],
        "available": int(parsed[2]) / KIB_GB,
        "total": int(parsed[3]) / KIB_GB
    }


# Get the current IP of the Pi
def network():
    # Get all system adapters
    adapters = ifaddr.get_adapters()

    for adapter in adapters:
        # Get the wireless interface
        try:
            if adapter.name.index("wl") != 0:
                continue
        except ValueError:
            continue

        # Find the IPv4 address
        for ip in adapter.ips:
            if ip.is_IPv4:
                return ip.ip

    raise ValueError("No wireless interfaces found")


# Get total number of clients
def clients():
    # TODO: add retrieval of clients
    return 1
