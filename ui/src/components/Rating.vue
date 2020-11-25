<template>
  <div class="rating-column rating-col">
    <b-icon-chevron-up
      :animation="animations[Math.floor(Math.random() * 6)]"
      @click.prevent="callback('up', entity)"
      shift-v="-7"
      class="h4"
      :style="`color: ${entity.your_reaction == 1 ? 'green' : 'white'}`"
    >
    </b-icon-chevron-up>
    <span>{{ entity.rating }}</span>
    <b-icon-chevron-down
      :animation="animations[Math.floor(Math.random() * 6)]"
      @click.prevent="callback('down', entity)"
      shift-v="2"
      class="h4"
      :style="`color: ${entity.your_reaction == -1 ? 'red' : 'white'}`"
    ></b-icon-chevron-down>
  </div>
</template>
<script>
export default {
  props: {
    callback: {
      type: Function,
      required: true,
    },
    entity: {
      type: Object,
      required: true,
    },
    type: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      animations: [
        "cylon",
        "cylon-vertical",
        "fade",
        "spin",
        "spin-reverse",
        "throb",
      ],
    };
  },
};
</script>

<style lang="scss">
.rating-col {
  display: flex;
  align-items: center;
  flex-direction: column;
  text-align: left;
  margin-left: -100px;
}
</style>
q
