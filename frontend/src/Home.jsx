function Home() {

return (
    <section className="bg-background min-h-screen flex flex-col justify-center items-center px-6 text-center space-y-6">
      <h1 className="text-4xl md:text-5xl font-bold text-primary leading-tight">
        Monitor Your Services. Catch Issues Before They Catch You.
      </h1>
      <p className="text-textSecondary text-lg md:text-xl max-w-xl">
        Uplytics keeps your websites, APIs, and apps running smoothly with real-time alerts and a public status page.
      </p>
      <button className="bg-orange-500 hover:bg-orange-600 text-white font-bold px-8 py-4 rounded-lg shadow-lg transition transform hover:-translate-y-1 hover:shadow-xl">
        Get Started - Free
      </button>
    </section>
  )
}



export default  Home
