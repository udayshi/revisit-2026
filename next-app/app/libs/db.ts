import sqlite3 from "sqlite3";
import { open } from "sqlite";

export async function openDB() {
  return open({
    filename: "users.db",
    driver: sqlite3.Database,
  });
}

// Initialize DB schema (run once)
export async function initDB() {
  const db = await openDB();
  await db.exec(`
    CREATE TABLE IF NOT EXISTS users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      email TEXT UNIQUE NOT NULL,
      username TEXT UNIQUE NOT NULL
    )
  `);
  await db.close();
}
initDB()