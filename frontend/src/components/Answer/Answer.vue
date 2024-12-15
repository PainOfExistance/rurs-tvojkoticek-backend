<template>
  <div>
    <h1 class="text-center my-4">Answer the Question</h1>

    <div class="form-group w-50 mx-auto p-3 border rounded shadow-sm bg-light">
      <form>
        <label class="form-label fw-bold" style="font-size: 0.9rem;">Author:</label>
        <input id="authorInput" class="form-control mb-2" placeholder="Enter your name" required style="font-size: 0.9rem;">

        <label class="form-label fw-bold" style="font-size: 0.9rem;">Your Answer:</label>
        <textarea id="answerInput" type="text" class="form-control mb-2" rows="5" placeholder="Write your answer here" required style="font-size: 0.9rem;"></textarea>

        <input type="submit" value="Submit" class="btn btn-secondary mt-2 d-block mx-auto" style="font-size: 0.9rem;">
      </form>
    </div>
  </div>
</template>

<script>
export default {

};
</script>

<style scoped>
.form-group {
  max-width: 500px;
  background-color: #f8f9fa;
  border-radius: 8px;
  font-size: 0.9rem;
}
</style>
