<template>
  <div :class="{ 'card-m': isMobile(), card: !isMobile() }">
    <b-skeleton
      type="input"
      variant="dark"
      animation="wave"
      width="100%"
      height="50px"
    ></b-skeleton>
    <b-skeleton variant="dark" animation="wave" width="0%"></b-skeleton>
    <b-row class="ml-1">
      <b-skeleton
        type="avatar"
        variant="dark"
        animation="wave"
        width="15px"
        height="15px"
        class="mr-1"
      ></b-skeleton>
      <b-skeleton variant="dark" animation="wave" width="15%"></b-skeleton>
    </b-row>
    <b-skeleton variant="dark" animation="wave" width="25%"></b-skeleton>
    <b-skeleton variant="dark" animation="wave" width="0%"></b-skeleton>
    <b-row class="ml-1">
      <b-skeleton
        type="avatar"
        variant="dark"
        animation="wave"
        width="15px"
        height="15px"
        class="mr-1"
      ></b-skeleton>
      <b-skeleton variant="dark" animation="wave" width="25%"></b-skeleton>
    </b-row>
    <b-skeleton variant="dark" animation="wave" width="20%"></b-skeleton>
    <b-skeleton variant="dark" animation="wave" width="0%"></b-skeleton>
    <b-row class="ml-1">
      <b-skeleton
        type="avatar"
        variant="dark"
        animation="wave"
        width="15px"
        height="15px"
        class="mr-1"
      ></b-skeleton>
      <b-skeleton variant="dark" animation="wave" width="10%"></b-skeleton>
    </b-row>
    <b-skeleton variant="dark" animation="wave" width="30%"></b-skeleton>
  </div>
</template>
<script>
export default {
  name: "CommentsSkeleton",
};
</script>
