// FurrDB CRUD Demo (Deno TypeScript)
// Demonstrates basic CRUD operations for a user profile

const HOST = "127.0.0.1";
const PORT = 7070;

// Helper to send a command and receive a response
async function sendCommand(conn: Deno.Conn, cmd: string): Promise<string> {
  const encoder = new TextEncoder();
  const decoder = new TextDecoder();
  await conn.write(encoder.encode(cmd + "\n"));
  let resp = "";
  while (!resp.endsWith("\n")) {
    const buf = new Uint8Array(1024);
    const n = await conn.read(buf);
    if (n === null) break;
    resp += decoder.decode(buf.subarray(0, n));
  }
  return resp.trim();
}

// Main demo logic
const conn = await Deno.connect({ hostname: HOST, port: PORT });
console.log("Connected to FurrDB\n");

// Create (SET user:1 name and email)
console.log("[CREATE] Set user:1 name and email");
console.log(await sendCommand(conn, 'SET user:1:name Alice'));
console.log(await sendCommand(conn, 'SET user:1:email alice@example.com'));

// Read (GET user:1 name and email)
console.log("\n[READ] Get user:1 name and email");
console.log('Name:', await sendCommand(conn, 'GET user:1:name'));
console.log('Email:', await sendCommand(conn, 'GET user:1:email'));

// Update (SET user:1 name)
console.log("\n[UPDATE] Update user:1 name");
console.log(await sendCommand(conn, 'SET user:1:name Alicia'));
console.log('Updated Name:', await sendCommand(conn, 'GET user:1:name'));

// Exists
console.log("\n[EXISTS] Check if user:1:email exists");
console.log('Exists:', await sendCommand(conn, 'EXISTS user:1:email'));

// Delete (DEL user:1 email)
console.log("\n[DELETE] Delete user:1:email");
console.log(await sendCommand(conn, 'DEL user:1:email'));
console.log('Email after delete:', await sendCommand(conn, 'GET user:1:email'));

// List all keys
console.log("\n[KEYS] List all keys");
console.log(await sendCommand(conn, 'KEYS'));

// Exit
await sendCommand(conn, 'EXIT');
conn.close();
console.log("\nConnection closed");

export {};