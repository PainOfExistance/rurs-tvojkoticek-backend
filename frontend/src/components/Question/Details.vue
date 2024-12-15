<template>
  <div v-if="loading"></div>
  <div v-else>
    <div class="post-data-container">
      <span>{{itemData.problem}}</span>
      <hr>
      <span>
        <span class="fw-bold">Avtor:</span> {{itemData.username}}
        <br>
        <span class="fw-bold">Datum objave:</span> {{new Date(itemData.date).toDateString()}}
      </span>
      <br>
      <div class="text-center mt-3">
        <button class="btn btn-danger" data-bs-toggle="modal" data-bs-target="#answerModal">Pomagaj</button>
      </div>
    </div>

    <div v-if="itemData.comments && itemData.comments !== []" class="mt-4 accordion" id="answerAccordion">
      <div class="accordion-item" v-for="(comment, index) in itemData.comments" :key="index">
        <h2 class="accordion-header">
          <button class="accordion-button" type="button" data-bs-toggle="collapse" :data-bs-target="'#' + index" :aria-expanded="index === 0" :aria-controls="index.toString()">
            <span class="mx-1"><span class="fw-bold">Avtor:</span> {{comment.username}},</span>
            <span class="mx-1">({{new Date(comment.date).toDateString()}})</span>
          </button>
        </h2>
        <div :id="index.toString()" class="accordion-collapse collapse show" data-bs-parent="#answerAccordion">
          <div class="accordion-body">
            {{comment.description}}
          </div>
        </div>

      </div>
    </div>
  </div>

  <div class="modal fade" id="answerModal" tabindex="-1" aria-labelledby="answerModalLabel" aria-hidden="true">
    <div class="modal-dialog">
      <div class="modal-content">
        <form @submit.prevent="postComment">
          <div class="modal-header">
            <h1 class="modal-title fs-5" id="answerModalLabel">Odgovor</h1>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">

              <label class="form-label">Avtor:</label>
              <input id="authorInput" class="form-control mb-2" placeholder="Vnesite vaše ime..." required style="font-size: 0.9rem;" v-model="formData.name">

              <label class="form-label">Vaše vprašanje:</label>
              <textarea id="storyInput" type="text" class="post-input-area form-control mb-2" rows="5" placeholder="Predlagajte odgovor" required style="font-size: 0.9rem;" v-model="formData.details"></textarea>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" id="closeModal" data-bs-dismiss="modal">Zapri</button>
            <button type="submit" class="btn btn-danger">Odgovori</button>
          </div>
        </form>
      </div>
    </div>
  </div>


</template>

<script>
import Cookies from 'js-cookie';
import Answers from "@/components/Answer/Answers.vue";
import axios from "axios";

export default {
  components: {Answers},
  props: ['id'],
  data() {
    return {
      itemData: null,
      loading: true,
      formData: {
        name: '',
        details: '',
      },
    };
  },
  methods: {
    checkLoginStatus() {
      this.isLoggedIn = !!Cookies.get('id');
    },
    getQuestionData() {
      if (this.isLoggedIn) {
        axios.get(`http://localhost:8080/post?post_id=${this.id}`).then(result => {
          if (result.status === 200) {
            this.itemData = result.data;
            this.loading = false;
            console.log(this.itemData);
          }
        }).catch(error => console.error(error));
      }
    },
    postComment() {
      let closeModalButton = document.getElementById('closeModal');
      axios.post('http://localhost:8080/comment', {
        post_id: this.id,
        username: this.formData.name,
        description: this.formData.details,
        date: new Date().toISOString().split('T')[0],
      }).then(response => {
        if (response.status === 200) {
          this.getQuestionData();
          closeModalButton.click();
        }
      }).catch(error => console.error(error));
    }
  },
  mounted() {
    this.checkLoginStatus();
    this.getQuestionData();
  },
};
</script>

<style scoped>
.btn-danger {
  background-color: #fbb2a3 !important;
  border-color: #fbb2a3 !important;
}

.btn-danger:hover {
  background-color: #fd9c8c !important;
  border-color: #fd9c8c !important;
}

.btn-secondary {
  background-color: #54627b;
}
</style>
