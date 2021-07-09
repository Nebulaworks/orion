import socket

def gethostname():
    return socket.gethostname()

def getlocaladdress():
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    # Doesnt matter what we try to connect to but just that we try
    s.connect(("8.8.8.8", 80))
    return s.getsockname()[0]
