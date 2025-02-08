import { Button } from "./ui/button";
import { ArrowRight, Globe } from "lucide-react";

export default function Hero() {
  return (
    <section className="container relative flex min-h-[calc(100vh-4rem)] max-w-screen-xl flex-col items-center justify-center space-y-10 py-24 text-center md:py-32">
      <div className="space-y-8">
        <div className="mx-auto flex max-w-fit items-center gap-2 overflow-hidden rounded-full border border-white/10 bg-white/5 px-4 py-2 backdrop-blur-lg">
          <Globe className="h-4 w-4 text-primary" />
          <p className="text-sm text-foreground/80">
            Supporting 100+ languages worldwide
          </p>
        </div>
        <h1 className="text-4xl font-bold tracking-tight sm:text-5xl md:text-6xl lg:text-7xl">
          Break Language Barriers
          <br />
          with <span className="text-gradient">Translatify</span>
        </h1>
        <p className="mx-auto max-w-[42rem] text-lg leading-relaxed text-foreground/70 sm:text-xl sm:leading-8">
          Seamlessly communicate with anyone, anywhere. Our AI-powered chat
          application translates messages in real-time, fostering global
          connections without language limitations.
        </p>
        <div className="flex justify-center gap-4">
          <Button
            size="lg"
            className="shadow-lg hover:shadow-xl transition-shadow"
          >
            Get Started
            <ArrowRight className="ml-2 h-4 w-4" />
          </Button>
          <Button
            size="lg"
            variant="outline"
            className="border-white/10 hover:bg-white/5"
          >
            Learn More
          </Button>
        </div>
      </div>
    </section>
  );
}
