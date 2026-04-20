import React from "react";

type LegalSection = {
  title: string;
  paragraphs: string[];
};

export default function LegalPage({
  title,
  updatedAt,
  intro,
  sections,
}: {
  title: string;
  updatedAt: string;
  intro: string;
  sections: LegalSection[];
}) {
  return (
    <div className="container mx-auto px-4 py-12 max-w-4xl">
      <header className="space-y-3">
        <h1 className="text-3xl md:text-4xl font-bold tracking-tight">
          {title}
        </h1>
        <p className="text-sm text-muted-foreground">Terakhir diperbarui: {updatedAt}</p>
        <p className="text-base text-muted-foreground leading-relaxed">{intro}</p>
      </header>

      <main className="mt-10 space-y-8">
        {sections.map((section, idx) => (
          <section key={`${idx}-${section.title}`} className="space-y-3">
            <h2 className="text-xl md:text-2xl font-semibold">{section.title}</h2>
            <div className="space-y-3">
              {section.paragraphs.map((p, pIdx) => (
                <p key={pIdx} className="text-sm md:text-base text-muted-foreground leading-relaxed">
                  {p}
                </p>
              ))}
            </div>
          </section>
        ))}
      </main>
    </div>
  );
}
