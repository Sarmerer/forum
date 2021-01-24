<template>
  <span>
    <b-form-group>
      <b-form-textarea
        aria-describedby="title-help-block"
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
      <b-form-text id="title-help-block" v-if="form.title">
        <span :style="`color: ${properTitleLength ? 'white' : 'red'}`"
          >{{ titleLength }}/{{ maxTitleLength }}
        </span>
      </b-form-text>
    </b-form-group>
    <b-form-group>
      <b-form-textarea
        aria-describedby="content-help-block"
        v-model="form.content"
        placeholder="* Cool content..."
        reqired
        :state="form.content ? properContentLength : null"
        rows="4"
        max-rows="50"
      ></b-form-textarea>
      <b-form-text id="content-help-block" v-if="form.content">
        <span :style="`color: ${properContentLength ? 'white' : 'red'}`"
          >{{ contentLength }}/{{ maxContentLength }}
        </span>
      </b-form-text>
    </b-form-group>
    <b-form-tags
      aria-describedby="tags-help-block"
      autocomplete="off"
      separator=" "
      remove-on-delete
      v-model="form.categories"
      tag-variant="dark"
      placeholder="Tags..."
      :tag-validator="tagValidator"
      @tag-state="onTagState"
    ></b-form-tags>
    <b-form-text id="tags-help-block">
      Lowercase, {{ minTagLength }}-{{ maxTagLength }} symbols
    </b-form-text>
  </span>
</template>
<script>
export default {
  props: { form: Object },
  watch: {
    validForm: function(val) {
      this.$emit("valid-form", val);
    },
  },
  computed: {
    validForm() {
      return (
        this.properTitleLength &&
        this.properContentLength &&
        !this.invalidTags.length &&
        !this.duplicateTags.length
      );
    },
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
  data() {
    return {
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
  created() {
    this.$emit("valid-form", this.validForm);
  },
  methods: {
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
  },
};
</script>
