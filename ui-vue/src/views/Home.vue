<template>
  <div class="home ">
    <div class="hero-image">
      <div class="hero-text">
        <h1>
          WELCOME<br />
          TO <span class="primary">FORUM</span>
        </h1>
      </div>
    </div>
    <div class="main-content">
      <div class="columns">
        <div class="post-col">
          <Posts />
        </div>
        <div class="info-col">
          <div class="card">
            <h3 class="primary">RECENT</h3>
            <hr />
            <p>col 2</p>
          </div>
          <div class="card">
            <h3 class="primary">CATEGORIES</h3>
            <hr />
            <Categories />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Posts from "@/components/Posts.vue";
import Categories from "@/components/Categories.vue";

export default {
  name: "Home",
  components: {
    Categories,
    Posts
  }
};
</script>
<style lang="scss">
/* The hero image */
.hero-image {
  background-image: url("../assets/img/home-hero.jpg");
  box-shadow: 0 5px 6px 2px rgba(10, 10, 10, 0.3);
  height: 25%;
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
  position: relative;
  margin-bottom: 25px;
}

.hero-text {
  text-align: center;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: large;
}

.columns {
  display: flex;
  flex-wrap: wrap;
  margin: 0 17%;
}

.columns > * {
  flex-basis: calc(calc(750px - 100%) * 999);
}

.card {
  margin: 20px;
  padding: 10px;
  background-color: rgba(255, 255, 255, 0.05);
  box-shadow: 5px 5px 6px 2px rgba(10, 10, 10, 0.3);
}

.post-col {
  flex-grow: 1;
}

.info-col {
  flex-grow: 0.4;
}

hr {
  opacity: 0.3;
}
</style>
