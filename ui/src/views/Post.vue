<template>
  <div>
    <div v-if="$route.params.id === 'new'" class="grid">
      <NewPost />
    </div>
    <div v-else class="grid">
      <ExistingPost :postID="postID" />
    </div>
  </div>
</template>

<script>
import ExistingPost from "@/components/ExistingPost";
import NewPost from "@/components/NewPost";
export default {
  data() {
    return {
      postID: Number.parseInt(this.$route.params.id)
    };
  },
  components: {
    ExistingPost,
    NewPost
  }
};
</script>

