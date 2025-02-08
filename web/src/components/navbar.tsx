import Link from "next/link";
import { Button } from "./ui/button";
import { Github } from "lucide-react";

export default function Navbar() {
  return (
    <header className="sticky top-0 z-50 w-full border-b border-white/10 bg-background/50 backdrop-blur-lg">
      <div className="container flex h-16 max-w-screen-2xl items-center">
        <nav className="flex flex-1 items-center space-x-6 text-sm font-medium">
          <Link href="/" className="mr-6 flex items-center space-x-2">
            <span className="font-bold text-2xl text-gradient">
              Translatify
            </span>
          </Link>
        </nav>
        <div className="flex items-center space-x-4">
          <Link
            href="https://github.com/HJyup/translatify"
            target="_blank"
            rel="noreferrer"
          >
            <Button
              variant="ghost"
              size="icon"
              className="text-foreground/70 hover:text-foreground"
            >
              <Github className="h-5 w-5" />
              <span className="sr-only">GitHub</span>
            </Button>
          </Link>
        </div>
      </div>
    </header>
  );
}
