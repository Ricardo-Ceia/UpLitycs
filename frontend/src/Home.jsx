import "./App.css";

function Home() {
  return (
    <section className="bg-background min-h-screen flex flex-col justify-center items-center px-6 text-center space-y-6">
      <ShinyText text="Monitor Your Services. Catch Issues Before They Catch You." />
      <p className="text-textSecondary text-lg md:text-xl max-w-xl">
        Uplytics keeps your websites, APIs, and apps running smoothly with real-time alerts and a public status page.
      </p>
      <button className="retro-button">
        Get Started - Free
      </button>
    </section>
  );
}

function ShinyText({ text }) {
  return (
    <h1 className="text-4xl md:text-5xl font-bold text-primary leading-tight">
      {text.split("").map((c, index) => (
        <span
          key={index}
          style={{ animationDelay: `${index * 0.1}s` }}
          className="shiny-letter"
        >
          {c}
        </span>
      ))}
    </h1>
  );
}

export default Home;

