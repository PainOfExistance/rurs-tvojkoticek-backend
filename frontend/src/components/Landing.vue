<template>
  <div class="landing-page">
    <!-- Decorative Background Shapes -->
    <div class="shape shape1"></div>
    <div class="shape shape2"></div>
    <div class="shape shape3"></div>

    <div class="content-wrapper">
      <img class="app-logo" src="@/assets/imgs/logo.png" alt="Logo" />
      <p class="subtitle">
        Dobrodošli v vašem toplem in zaupnem kotičku na spletu, kjer lahko odkrito delite svoje misli,
        poiščete razumevanje in podporo ter odkrijete, da v svojih izzivih niste sami.
      </p>
      <div class="buttons">
        <router-link to="/register" class="action-button register-button">Registracija</router-link>
        <router-link to="/login" class="action-button login-button">Prijava</router-link>
      </div>
      <!-- Optional: Add a short testimonial or quote -->
      <div class="testimonial">
        <blockquote>
          "TvojKoticek.si mi je pomagal najti notranji mir in se povezati s čudovitimi ljudmi."
        </blockquote>
        <cite>- Marija, uporabnica</cite>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "Landing",
};
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;600&display=swap');

html, body {
  margin: 0;
  padding: 0;
  height: 100%;
  width: 100%;
  overflow: hidden;
}

.landing-page {
  position: relative;
  width: 100vw;
  height: 100vh;
  background: linear-gradient(135deg, #fbcfe8, #fde2e4, #fbd4e7, #fbe9e7);
  background-size: 400% 400%;
  animation: gradientBG 20s ease infinite;
  display: flex;
  justify-content: center;
  align-items: center;
  font-family: 'Poppins', sans-serif;
  color: #444;
  text-align: center;
  overflow: hidden;
}

/* Decorative Shapes */
.shape {
  position: absolute;
  border-radius: 50%;
  opacity: 0.3;
  animation: float 25s infinite ease-in-out;
}

.shape1 {
  width: 200px;
  height: 200px;
  background: #fbd4e7;
  top: -50px;
  left: -50px;
}

.shape2 {
  width: 300px;
  height: 300px;
  background: #fde2e4;
  bottom: -100px;
  right: -100px;
}

.shape3 {
  width: 150px;
  height: 150px;
  background: #fbe9e7;
  top: 20%;
  right: 20%;
}

/* Content Wrapper */
.content-wrapper {
  position: relative;
  z-index: 1;
  max-width: 800px;
  width: 90%;
  padding: 40px 30px;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  animation: fadeIn 2s ease forwards;
}

/* Logo */
.app-logo {
  width: 150px;
  margin-bottom: 30px;
  animation: slideIn 1.5s ease forwards;
}

/* Subtitle */
.subtitle {
  font-size: 1.4rem;
  font-weight: 300;
  line-height: 1.8;
  margin-bottom: 50px;
  color: #333;
  animation: fadeInUp 2s ease forwards;
  opacity: 0;
}

/* Buttons */
.buttons {
  display: flex;
  gap: 30px;
  justify-content: center;
  margin-bottom: 30px;
}

.action-button {
  padding: 15px 40px;
  border-radius: 30px;
  text-decoration: none;
  font-weight: 500;
  font-size: 1.15rem;
  transition: background 0.3s ease, transform 0.3s ease;
  display: inline-block;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
  color: #fff;
}

.register-button {
  background-color: #8da4b9; /* Soft pastel blue-grey */
}

.register-button:hover {
  background-color: #7b93a9;
  transform: translateY(-3px);
}

.login-button {
  background-color: #f6a8c3; /* Soft pastel pink */
}

.login-button:hover {
  background-color: #ee94b3;
  transform: translateY(-3px);
}

/* Testimonial */
.testimonial {
  margin-top: 40px;
  font-style: italic;
  color: #555;
  animation: fadeInUp 2.5s ease forwards;
  opacity: 0;
}

.testimonial cite {
  display: block;
  margin-top: 10px;
  font-style: normal;
  color: #777;
}

/* Keyframes */
@keyframes gradientBG {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(20px); }
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes slideIn {
  from { opacity: 0; transform: translateX(-50px); }
  to { opacity: 1; transform: translateX(0); }
}

/* Responsive Adjustments */
@media (max-width: 768px) {
  .buttons {
    flex-direction: column;
    gap: 15px;
  }

  .action-button {
    width: 100%;
    text-align: center;
  }
}
</style>
