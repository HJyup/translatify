"use client";

import { ArrowRight } from "lucide-react";
import { buttonVariants } from "./ui/button";

const LoginButton = () => {
  return (
    <a
      href="/api/auth/login"
      className={buttonVariants({
        size: "lg",
        className: "shadow-lg hover:shadow-xl transition-shadow"
      })}
    >
      Try It Now
      <ArrowRight className="ml-2 h-4 w-4" />
    </a>
  );
};

export default LoginButton;
