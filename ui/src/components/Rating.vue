<template>
  <div class="rating-column rating-col">
    <b-icon-chevron-up
      @click.prevent="callback('up', entity)"
      class="h4"
      :style="`color: ${entity.your_reaction == 1 ? 'green' : 'white'}`"
    >
    </b-icon-chevron-up>
    <span>{{ entity.rating }}</span>
    <b-icon-chevron-down
      @click.prevent="callback('down', entity)"
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
