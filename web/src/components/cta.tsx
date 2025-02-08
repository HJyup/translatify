import { Button } from "./ui/button";
import { ArrowRight } from "lucide-react";

export default function CTA() {
  return (
    <section className="border-t border-white/10 bg-gradient-to-b from-background/50 to-background">
      <div className="container flex flex-col items-center gap-6 py-24 text-center md:py-32">
        <h2 className="font-bold text-3xl leading-[1.1] sm:text-3xl md:text-5xl text-gradient">
          Ready to connect globally?
        </h2>
        <p className="max-w-[42rem] leading-normal text-foreground/70 sm:text-xl sm:leading-8">
          Join millions of users who are breaking language barriers and
          fostering international connections with Translatify. Experience
          seamless communication across languages today.
        </p>
        <Button
          size="lg"
          className="mt-2 shadow-lg hover:shadow-xl transition-shadow"
        >
          Start Chatting Now
          <ArrowRight className="ml-2 h-4 w-4" />
        </Button>
      </div>
    </section>
  );
}
