import LoginButton from "./login-button";
import { Button } from "./ui/button";
import { Sparkles } from "lucide-react";

export default function Hero() {
  return (
    <section className="container relative flex min-h-[calc(100vh-4rem)] max-w-screen-xl flex-col items-center justify-center space-y-10 py-24 text-center md:py-32">
      <div className="space-y-8">
        <div className="mx-auto flex max-w-fit items-center gap-2 overflow-hidden rounded-full border border-white/10 bg-white/5 px-4 py-2 backdrop-blur-lg">
          <Sparkles className="h-4 w-4 text-primary" />
          <p className="text-sm text-foreground/80">
            Use AI to break language barriers
          </p>
        </div>
        <h1 className="text-3xl font-bold tracking-tight sm:text-4xl md:text-5xl lg:text-6xl">
          Chat Across Languages
          <br />
          with <span className="text-gradient">Translatify</span>
        </h1>
        <p className="mx-auto max-w-[42rem] text-lg leading-relaxed text-foreground/70 sm:text-xl sm:leading-8">
          Focus on the conversation, not the translation.
        </p>
        <div className="flex justify-center gap-4">
          <LoginButton />
          <Button
            size="lg"
            variant="outline"
            className="border-white/10 hover:bg-white/5"
            disabled
          >
            See How It Works [Coming Soon]
          </Button>
        </div>
      </div>
    </section>
  );
}
