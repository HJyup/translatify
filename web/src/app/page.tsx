'use client';

import Navbar from "../components/navbar";
import Hero from "../components/hero";
import Features from "../components/features";
import CTA from "../components/cta";
import { motion } from "framer-motion";

export default function Home() {
  return (
    <div className="relative min-h-screen bg-gradient-to-b from-background to-background/80">
      <div className="pointer-events-none fixed inset-0">
        <motion.div
          initial={{ scale: 0.8, opacity: 0 }}
          animate={{ 
            scale: [0.8, 1.1, 0.9],
            opacity: [0.3, 0.6, 0.3]
          }}
          transition={{
            duration: 8,
            repeat: Infinity,
            ease: "easeInOut"
          }}
          className="absolute right-0 top-0 -translate-y-1/4 translate-x-1/4 
          h-[600px] w-[600px] rounded-full bg-primary/15 blur-[120px]"
          aria-hidden="true"
        />
        <motion.div
          initial={{ scale: 0.9, opacity: 0 }}
          animate={{ 
            scale: [0.9, 1.2, 0.8],
            opacity: [0.3, 0.6, 0.3]
          }}
          transition={{
            duration: 8,
            repeat: Infinity,
            ease: "easeInOut",
            delay: 2
          }}
          className="absolute bottom-0 left-0 translate-y-1/4 -translate-x-1/4 
          h-[600px] w-[600px] rounded-full bg-blue-500/15 blur-[120px]"
          aria-hidden="true"
        />
      </div>

      <div className="relative z-10">
        <Navbar />
        <Hero />
        <Features />
        <CTA />
      </div>
    </div>
  );
}
