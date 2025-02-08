import { MessageSquare, Globe, Zap, Users } from "lucide-react";

const features = [
  {
    name: "Real-time Translation",
    description:
      "Instantly translate messages to and from over 100 languages with zero delay.",
    icon: Globe
  },
  {
    name: "Personal Chat Experience",
    description:
      "Enjoy intimate one-on-one conversations where you and your partner each communicate in your native language.",
    icon: Users
  },
  {
    name: "AI-Powered Accuracy",
    description:
      "Experience unmatched translation quality powered by state-of-the-art AI language models.",
    icon: Zap
  },
  {
    name: "Contextual Understanding",
    description:
      "Get translations that preserve meaning and tone across languages through advanced context analysis.",
    icon: MessageSquare
  }
];

export default function Features() {
  return (
    <section id="features" className="container space-y-16 py-24 md:py-32">
      <div className="mx-auto max-w-[58rem] text-center">
        <h2 className="font-bold text-3xl leading-[1.1] sm:text-3xl md:text-5xl text-gradient">
          Breaking Language Barriers
        </h2>
        <p className="mt-6 text-foreground/70 sm:text-lg">
          Experience seamless global communication with Translatify&apos;s
          powerful translation features.
        </p>
      </div>
      <div className="mx-auto grid max-w-5xl grid-cols-1 gap-8 md:grid-cols-2">
        {features.map(feature =>
          <div
            key={feature.name}
            className="glass-effect p-8 transition-all duration-300 hover:border-primary/50"
          >
            <div className="flex items-center gap-4">
              <div className="rounded-lg bg-primary/10 p-2">
                <feature.icon className="h-6 w-6 text-primary" />
              </div>
              <h3 className="font-bold text-lg text-foreground">
                {feature.name}
              </h3>
            </div>
            <p className="mt-4 text-foreground/70 leading-relaxed">
              {feature.description}
            </p>
          </div>
        )}
      </div>
    </section>
  );
}
