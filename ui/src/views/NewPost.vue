<template>
  <div class="grid">
    <div class="columns">
      <div v-if="isMobile()" class="info-col">
        <UserCard :userData="user" link />
      </div>
      <div
        :class="`main-col p-3 ${isMobile() ? 'card-m' : 'card'}`"
        id="new-post"
      >
        <h3 align="center">Create new post</h3>
        <b-overlay
          :show="requesting"
          rounded
          opacity="0.6"
          spinner-small
          variant="dark"
          spinner-variant="light"
          class="d-inline-block"
        >
          <template #overlay>
            <div class="text-center">
              <b-icon
                icon="stopwatch"
                font-scale="3"
                animation="cylon"
              ></b-icon>
              <p id="cancel-label">Please wait...</p>
            </div>
          </template>
          <small>* - required</small>
          <b-form @submit.prevent="submit">
            <PostForm :form="form" v-on:valid-form="validForm = $event" />
            <b-button
              :disabled="!validForm"
              type="submit"
              variant="info"
              class="mt-3"
              >Submit
            </b-button>
          </b-form>
        </b-overlay>
      </div>
      <div v-if="!isMobile()" class="info-col">
        <UserCard v-if="user" :userData="user" link />
      </div>
    </div>
  </div>
</template>
<script>
import PostForm from "@/components/forms/PostForm";
import UserCard from "@/components/UserCard";
import { mapGetters } from "vuex";
import api from "@/api/api";

export default {
  watch: {
    user(newVal) {
      if (newVal === null)
        this.$router.push({ name: "Auth", params: { prevRote: "/new-post" } });
    },
  },
  computed: {
    ...mapGetters({
      user: "auth/user",
    }),
  },
  components: { UserCard, PostForm },
  data() {
    return {
      validForm: false,
      requesting: false,
      form: {
        title: "",
        content: "",
        categories: [],
      },
    };
  },
  methods: {
    submit() {
      this.requesting = true;
      api
        .post("post/create", {
          title: this.form.title,
          content: this.form.content,
          categories: this.form.categories,
        })
        .then((response) => {
          this.$router.push({
            name: "Post",
            params: {
              id: response.data?.data?.id,
              postData: response?.data?.data,
            },
          });
        })
        .then(() => (this.requesting = false));
    },
  },
};
</script>
