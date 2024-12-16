<script>
import axios from "axios";
import JSZip from "jszip";

export default {
    name: "Video store",
    data() {
        return {
            videos: [],
        }
    },
    mounted() {
        this.runOnLoad();
    },
    methods: {
        async runOnLoad() {
            try {
                const response = await axios.get('http://localhost:8080/videostore/all', { responseType: 'arraybuffer' });
                const zip = await JSZip.loadAsync(response.data);
                const videoFiles = zip.file(/\.mp4$/);
                const metadataFiles = zip.file(/\.txt$/);

                for (let i = 0; i < metadataFiles.length; i++) {
                    const metadata = await metadataFiles[i].async("string");
                    const videoMetadata = this.parseMetadata(metadata);
                    const videoBlob = await videoFiles[i].async("blob");
                    const videoUrl = URL.createObjectURL(videoBlob);

                    this.videos.push({
                        ...videoMetadata,
                        url: videoUrl
                    });
                }

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
                    if(key == "Posted At")
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
        },
        goToDetails(id) {
            this.$router.push({ name: 'Video', params: { id } });
        }
    }
};
</script>

<template>
    <div>
        <div class="row">
            <div class="col-6 mb-4" v-for="(video, index) in videos" :key="index">
                <div class="card">
                    <div class="card-body video-card-container">
                        <h4 class="video-title">{{ video.video_name }}</h4>
                        <span class="card-text"><span class="fw-bold">Avtor:</span> {{ video.uploader }}</span>
                        <hr>
                        <video width="100%" controls>
                            <source :src="video.url" type="video/mp4">
                            Your browser does not support the video tag.
                        </video>
                        <br>
                        <span class="card-text"><span class="fw-bold">Datum objave:</span> {{ video.posted_at }}</span>
                        <br>
                        <span class="card-text"><span class="fw-bold">Opis:</span> {{ video.description }}</span>
                        <br>
                        <span class="card-text"><span class="fw-bold">Oznake:</span> {{ video.tags }}</span>
                        <div class="text-end">
                            <button class="btn btn-danger" @click="goToMore(video.ID)">Poglej več</button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style>
.video-card-container {
    border: 1px solid #54627b;
    border-radius: 5px;
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2);
}

.btn-danger {
    background-color: #fbb2a3 !important;
    border-color: #fbb2a3 !important;
}

.btn-danger:hover {
    background-color: #fd9c8c !important;
    border-color: #fd9c8c !important;
}
</style>
