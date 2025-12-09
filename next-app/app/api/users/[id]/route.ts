// app/api/users/[id]/route.ts
import { NextResponse } from "next/server";
import { openDB } from "@/app/libs/db";

export async function DELETE(
  request: Request,
  { params }: { 
    params: Promise<{ id: string }>  // ✅ Next.js 15 syntax
  }
) {
  const db = await openDB();
  try {
    const { id } = await params;  // ✅ Await the Promise
    const userId = Number(id);
    
    const result = await db.run("DELETE FROM users WHERE id = ?", userId);
    
    if (result.changes === 0) {
      return NextResponse.json({ error: "User not found" }, { status: 404 });
    }
    
    return NextResponse.json({ message: "User deleted successfully" });
  } finally {
    await db.close();
  }
}
