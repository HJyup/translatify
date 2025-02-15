"use client";

import { ArrowRight } from "lucide-react";
import { Button, buttonVariants } from "./ui/button";
import { useRouter } from "next/navigation";
const LoginButton = () => {
  const router = useRouter();

  return (
    <Button
      className={buttonVariants({
        size: "lg",
        className: "shadow-lg hover:shadow-xl transition-shadow"
      })}
      onClick={() => router.push("/login")}
    >
      Try It Now
      <ArrowRight className="ml-2 h-4 w-4" />
    </Button>
  );
};

export default LoginButton;
