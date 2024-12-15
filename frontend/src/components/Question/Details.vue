<template>
  <div v-if="loading" class="text-center mt-5">
    <div class="spinner-border text-primary" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
  </div>
  <div v-else>
    <!-- Post Details -->
    <div class="post-data-container p-4 rounded shadow-sm bg-white mb-4">
      <h4 class="fw-bold text-center text-secondary mb-3">{{ itemData.problem }}</h4>
      <hr />
      <div class="post-details text-muted">
        <p><span class="fw-bold">Avtor:</span> {{ itemData.username }}</p>
        <p><span class="fw-bold">Datum objave:</span> {{ new Date(itemData.date).toLocaleDateString("sl-SI") }}</p>
      </div>
      <div class="text-center mt-3">
        <button class="btn btn-danger px-4" data-bs-toggle="modal" data-bs-target="#answerModal">Pomagaj</button>
      </div>
    </div>

    <!-- Comments Accordion -->
    <div v-if="itemData.comments && itemData.comments.length" class="mt-4">
      <h5 class="fw-bold text-secondary mb-3">Komentarji</h5>
      <div class="accordion" id="answerAccordion">
        <div
          class="accordion-item rounded"
          v-for="(comment, index) in itemData.comments"
          :key="index"
        >
          <h2 class="accordion-header">
            <button
              class="accordion-button"
              type="button"
              data-bs-toggle="collapse"
              :data-bs-target="'#comment' + index"
              :aria-expanded="index === 0"
              :aria-controls="'comment' + index"
            >
              <span class="me-2"><strong>Avtor:</strong> {{ comment.username }}</span>
              <span>({{ new Date(comment.date).toLocaleDateString("sl-SI") }})</span>
            </button>
          </h2>
          <div
            :id="'comment' + index"
            class="accordion-collapse collapse"
            :class="{ show: index === 0 }"
            data-bs-parent="#answerAccordion"
          >
            <div class="accordion-body text-muted">
              {{ comment.description }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Modal for Comment -->
  <div
    class="modal fade"
    id="answerModal"
    tabindex="-1"
    aria-labelledby="answerModalLabel"
    aria-hidden="true"
  >
    <div class="modal-dialog">
      <div class="modal-content">
        <form @submit.prevent="postComment">
          <div class="modal-header">
            <h5 class="modal-title" id="answerModalLabel">Dodaj odgovor</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <label class="form-label">Avtor:</label>
            <input
              id="authorInput"
              class="form-control mb-2"
              placeholder="Vnesite vaše ime..."
              required
              v-model="formData.name"
            />

            <label class="form-label">Vaš odgovor:</label>
            <textarea
              id="storyInput"
              class="form-control mb-2"
              rows="5"
              placeholder="Vnesite vaš odgovor..."
              required
              v-model="formData.details"
            ></textarea>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Zapri</button>
            <button type="submit" class="btn btn-danger">Oddaj odgovor</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import Cookies from "js-cookie";
import axios from "axios";

export default {
  props: ["id"],
  data() {
    return {
      itemData: null,
      loading: true,
      formData: {
        name: "",
        details: "",
      },
    };
  },
  methods: {
    checkLoginStatus() {
      this.isLoggedIn = !!Cookies.get("id");
    },
    getQuestionData() {
      axios
        .get(`http://localhost:8080/post?post_id=${this.id}`)
        .then((result) => {
          if (result.status === 200) {
            this.itemData = result.data;
            this.loading = false;
          }
        })
        .catch((error) => console.error(error));
    },
    postComment() {
      axios
        .post("http://localhost:8080/comment", {
          post_id: this.id,
          username: this.formData.name,
          description: this.formData.details,
          date: new Date().toISOString().split("T")[0],
        })
        .then((response) => {
          if (response.status === 200) {
            this.getQuestionData();
            document.getElementById("closeModal").click();
          }
        })
        .catch((error) => console.error(error));
    },
  },
  mounted() {
    this.checkLoginStatus();
    this.getQuestionData();
  },
};
</script>

<style scoped>
.post-data-container {
  background-color: #ffffff;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-radius: 10px;
}

.accordion-button {
  background-color: #f8f9fa;
  font-size: 1rem;
}

.accordion-button:not(.collapsed) {
  background-color: #e9ecef;
  color: #495057;
}

.accordion-body {
  font-size: 0.95rem;
}

.btn-danger {
  background-color: #fbb2a3;
  border: none;
}

.btn-danger:hover {
  background-color: #fd9c8c;
}

.btn-secondary {
  background-color: #54627b;
  border: none;
}

.spinner-border {
  width: 3rem;
  height: 3rem;
}
</style>
