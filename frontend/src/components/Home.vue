<!-- src/components/Home.vue -->
<template>
  <div class="home-page">
    <!-- Dekorativne Oblike (same as Ask.vue) -->
    <div class="shape shape1"></div>
    <div class="shape shape2"></div>
    <div class="shape shape3"></div>

    <div class="content-wrapper">
      <img class="app-logo" src="@/assets/imgs/logo.png" alt="Logo" />
      <h2 class="form-title">Dobrodošli na TvojKoticek.si</h2>

      <div class="description-section">
        <h3>Vašem varnem prostoru za deljenje misli in iskanje podpore</h3>
        <p>
          Vsi se kdaj znajdemo v trenutkih, ko nas življenje preizkuša in potrebujemo prostor, kjer lahko brez obsojanja izrazimo svoje misli, težave in stiske. TvojKoticek.si je varen kotiček, kjer se lahko odprete, napišete svoje izkušnje in poiščete nasvet, tolažbo ali preprosto podporo. Tukaj lahko iskreno delite vse, kar vas teži, ter pridobite vpogled in podporo od drugih, ki so že morda doživeli podobne izkušnje.

        </p>
        <br />
        <p>
          Naša skupnost verjame v moč odprtega in razumevajočega dialoga. Ko napišete svoje misli, vam drugi uporabniki lahko komentirajo in s tem pomagajo videti stvari iz drugačnih perspektiv ali pa preprosto nudijo toplo besedo, ki včasih pomeni največ.
        </p>
        <hr />
        <p>
          Spletna stran TvojKoticek.si je namenjena vsem, ki si želijo:
          <span>
            - Podeliti svoje občutke in misli v varnem in podpornem okolju.
          </span>
          <span>
            - Pridobiti mnenje ali nasvet od ljudi, ki razumejo ali so šli skozi podobne izzive.
          </span>
          <span>
            - Zgraditi skupnost, ki temelji na spoštovanju, razumevanju in empatiji.
          </span>
        </p>
        <hr />
        <p>
          Vaša anonimnost je pri nas zagotovljena, saj želimo, da se počutite sproščeno in brez skrbi ob izražanju svojih občutkov.
        </p>
        <br />
        <p>
          Pridružite se nam na TvojKoticek.si – vašem prostoru za deljenje, podporo in korak naprej k notranjemu miru.        </p>
      </div>

      <div class="action-section mt-5">
        <router-link to="/ask" class="btn btn-primary">
          Vprašaj
        </router-link>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "Home",
};
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;600&display=swap');

.home-page {
  position: relative;
  width: 100%;
  min-height: calc(100vh - 80px); 
  background: linear-gradient(135deg);
  background-size: 400% 400%;
  animation: gradientBG 20s ease infinite;

  display: flex;
  flex-direction: column;
  align-items: center;
  font-family: 'Poppins', sans-serif;
  color: #444;
  text-align: center;
  padding: 20px;
  overflow-y: auto;
  box-sizing: border-box;
}

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
  top: 10px;
  left: 10px;
}

.shape2 {
  width: 300px;
  height: 300px;
  background: #fde2e4;
  bottom: 10px;
  right: 10px;
}

.shape3 {
  width: 150px;
  height: 150px;
  background: #fbe9e7;
  top: 20%;
  right: 20%;
}

.content-wrapper {
  position: relative;
  z-index: 1;
  max-width: 800px;
  width: 90%;
  padding: 40px 30px;
  background: rgba(241, 241, 241, 0.473); 
  backdrop-filter: blur(10px);
  border-radius: 20px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);

  display: flex;
  flex-direction: column;
  align-items: center; 
  margin-bottom: 40px;
}

.app-logo {
  width: 150px;
  margin-bottom: 20px;
}

.form-title {
  font-size: 2rem;
  font-weight: 600;
  margin-bottom: 30px;
  color: #333;
}

.description-section {
  text-align: left;
  color: #333;
  font-size: 1rem;
  max-height: none;
  overflow-y: visible;
}

.description-section h3 {
  color: #333;
  margin-bottom: 15px;
}

.action-section {
  margin-top: 30px;
  width: 100%;
  display: flex;
  justify-content: center;
}

.btn-primary {
  background-color: #8da4b9;
  border: none;
  color: white;
  padding: 15px 30px;
  border-radius: 30px;
  font-weight: 600;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.3s ease, transform 0.3s ease;
  text-decoration: none;
  display: inline-block;
  margin: 0 10px;
}
.btn-primary:hover {
  background-color: #7b93a9;
  transform: translateY(-3px);
}

@keyframes gradientBG {
  0%   { background-position: 0% 50%; }
  50%  { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50%      { transform: translateY(20px); }
}

@media (max-width: 768px) {
  .content-wrapper {
    padding: 30px 20px;
    max-width: 90%;
    height: auto;
  }
  .form-title {
    font-size: 1.75rem;
  }
  .btn-primary {
    width: 100%;
    margin-bottom: 15px;
  }
}
</style>
