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
          @hidden="onHidden"
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
          <b-form @submit.prevent="onSubmit">
            <b-form-group label-for="title">
              <small>* - required</small>
              <b-form-textarea
                class="mt-1"
                v-model="form.title"
                autocomplete="off"
                rows="1"
                :state="form.title ? properTitleLength : null"
                max-rows="8"
                no-resize
                required
                placeholder="* Catchy title..."
              ></b-form-textarea>
              <small
                v-if="form.title"
                :style="`color: ${properTitleLength ? 'green' : 'red'}`"
                >{{ titleLength }}/{{ maxTitleLength }}
              </small>
            </b-form-group>
            <b-form-group id="input-group-2" label-for="input-2" fluid>
              <b-form-textarea
                id="textarea-auto-height"
                v-model="form.content"
                placeholder="* Cool content..."
                reqired
                :state="form.content ? properContentLength : null"
                rows="4"
                max-rows="50"
              ></b-form-textarea>
              <small
                v-if="form.content"
                :style="`color: ${properContentLength ? 'green' : 'red'}`"
                >{{ contentLength }}/{{ maxContentLength }}
              </small>
            </b-form-group>
            <b-form-tags
              autocomplete="off"
              remove-on-delete
              v-model="form.categories"
              tag-variant="dark"
              :placeholder="
                `Lowercase, ${minTagLength}-${maxTagLength} symbols`
              "
              :tag-validator="tagValidator"
              @tag-state="onTagState"
            ></b-form-tags>

            <b-button
              :disabled="
                !properTitleLength ||
                  !properContentLength ||
                  invalidTags.length > 0 ||
                  duplicateTags.length > 0
              "
              type="submit"
              variant="info"
              class="mt-3"
              >Submit</b-button
            >
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
import api from "@/router/api";
import UserCard from "@/components/UserCard";
import { mapGetters } from "vuex";

export default {
  computed: {
    ...mapGetters({
      user: "auth/user",
    }),
    titleLength() {
      return this.form.title.replace(/(\r\n|\n|\r|\s)/g, "").length;
    },
    contentLength() {
      return this.form.content.replace(/(\r\n|\n|\r|\s)/g, "").length;
    },
    properTitleLength() {
      return (
        this.titleLength >= this.minTitleLength &&
        this.titleLength <= this.maxTitleLength
      );
    },
    properContentLength() {
      return (
        this.contentLength >= this.minContentLength &&
        this.contentLength <= this.maxContentLength
      );
    },
  },
  beforeRouteLeave(to, from, next) {
    if (
      (this.form.title.length ||
        this.form.content.length ||
        this.form.categories.length) &&
      this.user
    ) {
      if (window.confirm("Are tou sure?")) next();
    } else {
      next();
    }
  },
  components: { UserCard },
  data() {
    return {
      form: {
        title: "",
        amount: 1,
        content: "",
        categories: [],
      },

      requesting: false,

      invalidTags: [],
      duplicateTags: [],

      minTitleLength: 5,
      maxTitleLength: 300,

      minContentLength: 5,
      maxContentLength: 2000,

      minTagLength: 3,
      maxTagLength: 20,
    };
  },
  methods: {
    onSubmit() {
      this.requesting = true;
      api
        .post("post/create", {
          title: this.form.title,
          content: this.form.content,
          categories: this.form.categories,
        })
        .then((response) => {
          this.resetForm();
          this.$router.push({
            name: "Post",
            params: { id: response.data.data.id, postData: response.data.data },
          });
        })
        .catch(() => {
          this.$router.push("/");
        })
        .then(() => (this.requesting = false));
    },
    onTagState(_valid, invalid, duplicate) {
      this.invalidTags = invalid;
      this.duplicateTags = duplicate;
    },
    tagValidator(tag) {
      return (
        tag === tag.toLowerCase() &&
        tag.length >= this.minTagLength &&
        tag.length <= this.maxTagLength
      );
    },
    resetForm() {
      this.form.title = "";
      this.form.content = "";
      this.form.categories = [];
    },
  },
};
</script>