"use client";

import  { useState, useEffect, useLayoutEffect, useCallback, useMemo } from "react";

interface Props {
  initialCount?: number;
}

export default function LifeCycleDemo({ initialCount = 0 }: Props) {
  const [count, setCount] = useState(initialCount);

  // 1. Mounting (replaces constructor + componentDidMount)
  useEffect(() => {
    console.log("Mounting: componentDidMount equivalent");
    return () => {
      console.log("Unmounting: componentWillUnmount equivalent");
    };
  }, []); // empty deps = mount/unmount only

  // 2. Props change (getDerivedStateFromProps equivalent)
  useEffect(() => {
    console.log("Props changed: getDerivedStateFromProps equivalent");
  }, [initialCount]);

  // 3. Update after render (componentDidUpdate equivalent)
  useEffect(() => {
    console.log("Updated: componentDidUpdate equivalent");
  }, [count]);

  // 4. Sync DOM read/write (getSnapshotBeforeUpdate equivalent)
  useLayoutEffect(() => {
    console.log("Layout effect: getSnapshotBeforeUpdate equivalent");
    console.log("Current count value:", count);
  }, [count]);

  // 5. Prevent unnecessary re-renders (shouldComponentUpdate equivalent)
  const increment = useCallback(() => {
    setCount((prev) => prev + 1);
  }, []);

  // 6. Memoized value (render optimization)
  const displayText = useMemo(() => `Count: ${count}`, [count]);

  console.log("render");

  return (
    <div>
      <h2>Lifecycle Demo (Hooks)</h2>
      <p>{displayText}</p>
      <button onClick={increment}>Increment</button>
    </div>
  );
}
