import { NextResponse } from "next/server";
import { openDB } from "@/app/libs/db";

export async function GET() {
  const db = await openDB();
  const users = await db.all("SELECT * FROM users");
  await db.close();
  return NextResponse.json(users);
}


export async function POST(request: Request) {
  const { email, username } = await request.json();
  const db = await openDB();
  try {
    const result = await db.run(
      "INSERT INTO users (email, username) VALUES (?, ?)",
      email,
      username
    );
    const user = await db.get("SELECT * FROM users WHERE id = ?", result.lastID);
    await db.close();
    return NextResponse.json(user, { status: 201 });
  } catch (error: any) {
    await db.close();
    return NextResponse.json({ error: error.message }, { status: 400 });
  }
}



