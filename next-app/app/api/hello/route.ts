import { NextResponse } from "next/server";

export async function GET(request: Request) {
  return NextResponse.json({ message: "Hello from GET" });
}

interface PostRequestBody {
  message: string;
}
export async function POST(request: Request) {
  const body: PostRequestBody = await request.json();
  console.log("Received in API:", body.message);
  return NextResponse.json({ received: body }, { status: 201 });
}