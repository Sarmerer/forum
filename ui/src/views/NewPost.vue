<template>
  <div class="grid">
    <div class="columns">
      <div v-if="isMobile()" class="info-col">
        <UserCard :userData="user" link />
      </div>
      <div
        class="main-col p-3"
        :class="{ 'card-m': isMobile(), card: !isMobile() }"
        id="new-post"
      >
        <h4 align="center">Create new post</h4>
        <PostForm>
          <template slot="buttons" slot-scope="props">
            <b-button
              :disabled="!props.validForm"
              type="submit"
              variant="info"
              class="mt-3"
              >Submit
            </b-button>
          </template>
        </PostForm>
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
};
</script>
