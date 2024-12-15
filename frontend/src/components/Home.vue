<script setup>

</script>

<template>
  <div>
    <div class="description-section">
      <h3>Dobrodošli na TvojKoticek.si – Vašem varnem prostoru za deljenje misli in iskanje podpore</h3>
      <span>
        Vsi se kdaj znajdemo v trenutkih, ko nas življenje preizkuša in potrebujemo prostor, kjer lahko brez obsojanja izrazimo svoje misli, težave in stiske. TvojKoticek.si je varen kotiček, kjer se lahko odprete, napišete svoje izkušnje in poiščete nasvet, tolažbo ali preprosto podporo. Tukaj lahko iskreno delite vse, kar vas teži, ter pridobite vpogled in podporo od drugih, ki so že morda doživeli podobne izkušnje.
      </span>
      <br>
      <br>
      <span>
        Naša skupnost verjame v moč odprtega in razumevajočega dialoga. Ko napišete svoje misli, vam drugi uporabniki lahko komentirajo in s tem pomagajo videti stvari iz drugačnih perspektiv ali pa preprosto nudijo toplo besedo, ki včasih pomeni največ.
      </span>
      <hr>
      <span>
        Spletna stran TvojKoticek.si je namenjena vsem, ki si želijo:
        <br>
        - Podeliti svoje občutke in misli v varnem in podpornem okolju.
        <br>
        - Pridobiti mnenje ali nasvet od ljudi, ki razumejo ali so šli skozi podobne izzive.
        <br>
        - Zgraditi skupnost, ki temelji na spoštovanju, razumevanju in empatiji.
      </span>
      <hr>
      <span>Vaša anonimnost je pri nas zagotovljena, saj želimo, da se počutite sproščeno in brez skrbi ob izražanju svojih občutkov.</span>
      <br>
      <span>Pridružite se nam na TvojKoticek.si – vašem prostoru za deljenje, podporo in korak naprej k notranjemu miru.</span>
    </div>

    <div class="action-section mt-5">
      <div v-if="isLoggedIn" class="row row-cols-1">
        <div class="col text-center">
          <button class="btn btn-danger"><router-link class="button-link-text" to="/ask">Vprašaj</router-link></button>
        </div>
      </div>
      <div v-if="!isLoggedIn" class="row row-cols-2">
        <div class="col text-center">
          <button class="btn btn-secondary"><router-link class="button-link-text" to="/login">Prijava</router-link></button>
        </div>
        <div class="col text-center">
          <button class="btn btn-secondary"><router-link class="button-link-text" to="/register">Registracija</router-link></button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Cookies from 'js-cookie';

export default {
  data() {
    return {
      isLoggedIn: false,
    };
  },
  methods: {
    checkLoginStatus() {
      this.isLoggedIn = !!Cookies.get('id');
    }
  },
  mounted() {
    this.checkLoginStatus();
  },
};
</script>

<style scoped>
.button-link-text {
  color: white;
  text-decoration: none;
}

.btn-danger {
  background-color: #fbb2a3;
  border-color: #fbb2a3;
}

.btn-danger:hover {
  background-color: #fd9c8c;
  border-color: #fd9c8c;
}

.btn-secondary {
  background-color: #54627b;
}
</style>
