<script>
import { ref } from 'vue'
import axios from "axios";
import Cookies from "js-cookie";


export default {
  name: "Posts",
  data() {
    return {
      data: null,
    }
  },
  mounted() {
    this.runOnLoad();
  },
  methods: {
    runOnLoad() {
      axios.get('http://localhost:8080/posts').then(respone => {
        this.data = respone.data;
      }).catch(error => console.error(error));
    },
    goToDetails(id) {
      this.$router.push({name: 'Details', params: { id }});
    }
  }
};

</script>

<template>
  <div>
    <div class="row">
      <div class="col-6 mb-4" v-for="(item, index) in data" :key="index">
        <div class="card card-hover">
          <div class="card-body post-card-container">
            <h4 class="card-title">{{item.problem}}</h4>
            <hr>
            <span class="card-text"><span class="fw-bold">Datum objave:</span> {{item.date}}</span>
            <br>
            <span class="card-text"><span class="fw-bold">Avtor:</span> {{item.username}}</span>
            <div class="text-end">
              <button class="btn btn-danger" @click="goToDetails(item.ID)">Pomagaj</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>


</template>

<style>
  .post-card-container {
    border: 1px solid #54627b;
    border-radius: 5px;
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2);
  }

  .btn-danger {
    background-color: #fbb2a3 !important;
    border-color: #fbb2a3 !important;
  }

  .btn-danger:hover {
    background-color: #fd9c8c !important;
    border-color: #fd9c8c !important;
  }

  .card-hover:hover {
    background: rgba(162, 190, 239, 0.54);
    transition: background-color 0.5s ease;
  }
</style>
