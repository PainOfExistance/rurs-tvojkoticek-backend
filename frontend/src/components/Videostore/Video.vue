<script>
import axios from "axios";
import JSZip from "jszip";

export default {
    name: "Video",
    data() {
        return {
            video: null
        }
    },
    mounted() {
        this.fetchVideo();
    },
    methods: {
        async fetchVideo() {
            try {
                const videoId = this.$route.params.id;
                const response = await axios.get(`http://localhost:8080/videostore/video${videoId}`, { responseType: 'arraybuffer' });
                const zip = await JSZip.loadAsync(response.data);
                const videoFile = zip.file(/\.mp4$/)[0];
                const metadataFile = zip.file(/\.txt$/)[0];

                const metadata = await metadataFile.async("string");
                const videoMetadata = this.parseMetadata(metadata);
                const videoBlob = await videoFile.async("blob");
                const videoUrl = URL.createObjectURL(videoBlob);

                this.video = {
                    ...videoMetadata,
                    url: videoUrl
                };
            } catch (error) {
                console.error(error);
            }
        },
        parseMetadata(metadata) {
            const lines = metadata.split('\n');
            const videoMetadata = {};

            lines.forEach(line => {
                const [key, value] = line.split(': ');
                if (key && value) {
                    if (key == "Posted At")
                        videoMetadata[key.toLowerCase().replace(' ', '_')] = this.parseDate(value);
                    else
                        videoMetadata[key.toLowerCase().replace(' ', '_')] = value;
                }
            });

            return videoMetadata;
        },
        parseDate(dateString) {
            const date = new Date(dateString);
            const day = date.getDate().toString().padStart(2, '0');
            const month = (date.getMonth() + 1).toString().padStart(2, '0');
            const year = date.getFullYear();
            const hours = date.getHours().toString().padStart(2, '0');
            const minutes = date.getMinutes().toString().padStart(2, '0');
            return `${day}.${month}.${year}  ${hours}:${minutes}`;
        }
    }
};
</script>

<template>
    <div v-if="video" class="video-container">
        <h1 class="text-center">{{ video.video_name }}</h1>
        <video width="100%" controls>
            <source :src="video.url" type="video/mp4">
            Your browser does not support the video tag.
        </video>
        <p><strong>Avtor:</strong> {{ video.uploader }}</p>
        <p><strong>Datum objave:</strong> {{ video.posted_at }}</p>
        <p><strong>Opis:</strong> {{ video.description }}</p>
        <!--<p><strong>Oznake:</strong> {{ video.tags }}</p>-->
    </div>
</template>

<style>
.video-container {
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
}
</style>
