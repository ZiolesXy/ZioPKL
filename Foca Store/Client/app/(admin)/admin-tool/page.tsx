import React from "react";
import Image from "next/image";
import { Card } from "@/components/ui/card"; // Menggunakan base Card dari Shadcn

export default function AdminPage() {
  return (
    <div className="flex justify-center items-center bg-black py-20">
      {/* Container Utama (Card) */}
      <Card className="group relative flex h-[500px] w-[333px] items-end justify-center px-9 perspective-[2500px]">
        
        {/* Layer 1: Background Card */}
        <div 
          className="absolute inset-0 z-[-1] bg-cover bg-center transition-all duration-500 
                     group-hover:opacity-20 
                     group-hover:[transform:perspective(1000px)_translateY(-10%)_rotateX(45deg)_translateZ(0)]"
          style={{ backgroundImage: "url('/alucard.jpeg')" }}
        >
          {/* Title Layer */}
          <div className="absolute bottom-0 w-full bg-black p-5 text-center">
            <h3 className="text-2xl font-extrabold uppercase text-white tracking-widest">
              Alucard
            </h3>
          </div>
        </div>

        {/* Layer 2: Character Image (Pop out) */}
        <div className="absolute z-[-1] w-[150%] opacity-0 transition-all duration-500 
                        group-hover:opacity-100 
                        group-hover:[transform:translate3d(0,-30px,100px)]">
          <img 
            src="/alucard_bg.png" 
            alt="Character" 
            className="w-full object-contain"
          />
        </div>

      </Card>
    </div>
  );
}