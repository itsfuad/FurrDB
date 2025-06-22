import socket

HOST = '127.0.0.1'
PORT = 7070

def send_command(sock, cmd):
    sock.sendall((cmd + '\n').encode())
    resp = b''
    while not resp.endswith(b'\n'):
        chunk = sock.recv(4096)
        if not chunk:
            break
        resp += chunk
    return resp.decode().strip()

with socket.create_connection((HOST, PORT)) as s:
    print('Connected to FurrDB\n')

    # CREATE
    print('[CREATE] Set user:1 name and email')
    print(send_command(s, 'SET user:1:name Alice'))
    print(send_command(s, 'SET user:1:email alice@example.com'))

    # READ
    print('\n[READ] Get user:1 name and email')
    print('Name:', send_command(s, 'GET user:1:name'))
    print('Email:', send_command(s, 'GET user:1:email'))

    # UPDATE
    print('\n[UPDATE] Update user:1 name')
    print(send_command(s, 'SET user:1:name Alicia'))
    print('Updated Name:', send_command(s, 'GET user:1:name'))

    # EXISTS
    print('\n[EXISTS] Check if user:1:email exists')
    print('Exists:', send_command(s, 'EXISTS user:1:email'))

    # DELETE
    print('\n[DELETE] Delete user:1:email')
    print(send_command(s, 'DEL user:1:email'))
    print('Email after delete:', send_command(s, 'GET user:1:email'))

    # KEYS
    print('\n[KEYS] List all keys')
    print(send_command(s, 'KEYS'))

    # EXIT
    send_command(s, 'EXIT')
    print('\nConnection closed')