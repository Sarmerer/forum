<template>
  <div>
    <div class="columns">
      <div class="info-col">
        <div class="card">
          <h3 class="primary">AUTHOR</h3>
          <p>Author info</p>
        </div>
      </div>
      <div class="post-col">
        <PostSection :postID="postID" :hasPermission="hasPermission" />
        <CommentsSection :postID="postID" :hasPermission="hasPermission" />
      </div>
    </div>
  </div>
</template>
<script>
import CommentsSection from "@/components/CommentsSection";
import PostSection from "@/components/PostSection";
import { mapGetters } from "vuex";

export default {
  components: {
    PostSection,
    CommentsSection,
  },
  computed: {
    ...mapGetters({
      user: "auth/user",
    }),
  },
  data() {
    return {
      postID: Number.parseInt(this.$route.params.id),
    };
  },
  methods: {
    hasPermission(author) {
      let self = this;
      return () => (self.user ? author == self.user.id || self.user.role > 0 : false);
    },
  },
};
</script>
