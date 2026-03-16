"use client";

import React, { useState, useEffect } from 'react';
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Sparkles, Music, Heart } from "lucide-react";

// Koleksi Quotes ala Kafka & Motivasi
const quotes = [
  "Listen... The world is quiet now. Take a deep breath.",
  "You don't need to know the future, you just need to walk toward it.",
  "Destiny is a fickle thing, but right now, you are exactly where you need to be.",
  "Close your eyes. Let the symphony of the stars calm your mind.",
  "You’ ve done enough for today. Even hunters need their rest.",
  "Patience is a virtue, especially when you're waiting for your own bloom."
];

export default function KafkaEncouragementPage() {
  const [quote, setQuote] = useState("");

  useEffect(() => {
    setQuote(quotes[Math.floor(Math.random() * quotes.length)]);
  }, []);

  const shuffleQuote = () => {
    const newQuote = quotes[Math.floor(Math.random() * quotes.length)];
    setQuote(newQuote);
  };

  return (
    <div className="min-h-screen bg-[#0f0a14] text-purple-100 font-sans selection:bg-purple-900">
      {/* Background Overlay Estetik */}
      <div className="fixed inset-0 bg-[url('https://images.alphacoders.com/131/1310620.jpeg')] bg-cover bg-fixed opacity-10 blur-sm pointer-events-none" />

      <main className="relative z-10 max-w-6xl mx-auto p-6 space-y-12">
        
        {/* Header Section */}
        <header className="text-center pt-10 space-y-4">
          <h1 className="text-5xl md:text-7xl font-bold tracking-tighter bg-gradient-to-r from-purple-400 to-pink-600 bg-clip-text text-transparent drop-shadow-sm">
            Kafka's Sanctuary
          </h1>
          <p className="text-muted-foreground italic text-lg">"Just listen to the sound of my voice..."</p>
        </header>

        {/* Main Feature: Quote Card */}
        <div className="flex justify-center">
          <Card className="w-full max-w-3xl bg-black/60 border-purple-900/50 backdrop-blur-md overflow-hidden transition-all hover:border-purple-500/50">
            <div className="md:flex">
              <div className="md:w-1/3 h-64 md:h-auto overflow-hidden">
                <img 
                  src="https://api.dicebear.com/7.x/identicon/svg?seed=Kafka" // Ganti dengan gambar portrait Kafka
                  alt="Kafka Portrait"
                  className="w-full h-full object-cover transition-transform duration-700 hover:scale-110"
                />
              </div>
              <CardContent className="md:w-2/3 p-8 flex flex-col justify-center space-y-6">
                <div className="space-y-2">
                  <Sparkles className="text-purple-400 w-6 h-6 mb-2" />
                  <p className="text-2xl font-medium leading-relaxed italic text-purple-50">
                    "{quote}"
                  </p>
                </div>
                <Button 
                  onClick={shuffleQuote}
                  className="w-fit bg-purple-800 hover:bg-purple-700 text-white rounded-full px-6 transition-all active:scale-95"
                >
                  Listen Again
                </Button>
              </CardContent>
            </div>
          </Card>
        </div>

        {/* Kafka Image Gallery / Cards Section */}
        <section className="space-y-6">
          <div className="flex items-center space-x-2">
            <Heart className="text-pink-500 fill-pink-500" />
            <h2 className="text-2xl font-semibold">Fragments of Memory</h2>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {[1, 2, 3, 4, 5, 6].map((i) => (
              <Card key={i} className="group overflow-hidden border-none bg-transparent">
                <div className="relative aspect-[4/5] overflow-hidden rounded-xl border-2 border-purple-900/30 group-hover:border-purple-500/50 transition-all">
                  <img 
                    src={"/mybini.jpeg"} // Ganti dengan URL gambar Kafka favoritmu
                    alt={`Kafka Art ${i}`}
                    className="object-cover w-full h-full grayscale-[20%] group-hover:grayscale-0 transition-all duration-500 group-hover:scale-105"
                  />
                  <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity p-4 flex flex-end items-end">
                    <p className="text-xs text-purple-300 uppercase tracking-widest font-bold">Stellaron Hunter</p>
                  </div>
                </div>
              </Card>
            ))}
          </div>
        </section>

        {/* Footer Peace Message */}
        <footer className="py-20 text-center space-y-4">
          <div className="flex justify-center gap-4 text-purple-400">
             <Music className="animate-bounce" />
          </div>
          <p className="max-w-md mx-auto text-sm text-muted-foreground leading-relaxed">
            Halaman ini didedikasikan untuk ketenanganmu. Seperti alunan biola Kafka, biarkan rasa lelahmu memudar sejenak.
          </p>
        </footer>
      </main>
    </div>
  );
}