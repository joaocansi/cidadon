"use client";

import { SendIcon } from "lucide-react";
import { useEffect, useState } from "react";

export function ProtocolTicker() {
  const [count, setCount] = useState(48213);

  useEffect(() => {
    const id = setInterval(() => {
      setCount((c) => c + Math.floor(Math.random() * 2) + 1);
    }, 4200);
    return () => clearInterval(id);
  }, []);

  return (
    <span className="protocol-ticker">
      <SendIcon className="icon mr-1.5 inline-block h-3.5 w-3.5" />
      Demandas registradas hoje no Brasil: <b>{count.toLocaleString("pt-BR")}</b>
    </span>
  );
}
