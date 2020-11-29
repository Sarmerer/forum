<template>
  <div :class="`user-card ${isMobile() ? 'card-m' : 'card'}`">
    <img :src="userData.avatar" alt="avatar" />
    <h3 class="primary">
      {{ userData.display_name }}
      <b-badge v-if="userData.role == 2" class="background-variant"
        >Admin</b-badge
      >
    </h3>
    <p>Last active: {{ userData.last_online }}</p>
    <p>Rating: {{ userData.rating }}</p>
    <p>Posts: {{ userData.posts }}</p>
    <p>Comments: {{ userData.comments }}</p>
  </div>
</template>
<script>
export default {
  name: "user-card",
  props: { userData: { type: Object } },
};
</script>
<style lang="scss" scoped>
.user-card {
  cursor: pointer;
}

.user-card img {
  display: block;
  margin-left: auto;
  margin-right: auto;
  margin-bottom: 15px;
  width: 150px;
  height: 150px;
  border-radius: 150px;
}
</style>
