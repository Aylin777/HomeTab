package pl.systemz.tasktab.api;

import android.content.Context;

import java.util.List;

import retrofit2.Call;
import retrofit2.Retrofit;
import retrofit2.converter.gson.GsonConverterFactory;
import retrofit2.http.GET;
import retrofit2.http.POST;
import retrofit2.http.Path;

public class Client {
    Context context;
    public static final String API_URL = "http://192.168.2.88:3000/api/v1/";
    private static Client instance = null;
    private GitHub github;

    public static class Contributor {
        public final String login;
        public final int contributions;

        public Contributor(String login, int contributions) {
            this.login = login;
            this.contributions = contributions;
        }
    }

    public static class Timer {
        public final int id;
        public final String name;
        public final List<String> tags;
        public final int seconds;
        public final boolean inProgress;

        public Timer(int id, String name, List<String> tags, int seconds, boolean inProgress) {
            this.id = id;
            this.name = name;
            this.tags = tags;
            this.seconds = seconds;
            this.inProgress = inProgress;
        }
    }

    public interface GitHub {
        @GET("/repos/{owner}/{repo}/contributors")
        Call<List<Contributor>> contributors(
                @Path("owner") String owner,
                @Path("repo") String repo);

        @GET("timer")
        Call<List<Timer>> timers();

        @GET("timer/{id}")
        Call<Timer> timerInfo(
                @Path("id") int id
        );

        @POST("timer/{id}/start")
        Call<Timer> timerStart(
                @Path("id") int id
        );

        @POST("timer/{id}/stop")
        Call<Timer> timerStop(
                @Path("id") int id
        );
    }

    private Client() {
        buildRetrofit(API_URL);
    }

    public static Client getInstance() {
        if (instance == null) {
            instance = new Client();
        }
        return instance;
    }

//    private SharedPreferences prefs;

//    private String authToken;
//
//    public ServiceInterceptor(String authToken) {
//        this.authToken = authToken;
//    }

//    class ServiceInterceptor implements Interceptor{
//        @NonNull
//        @Override
//        public okhttp3.Response intercept(@NonNull Chain chain) throws IOException {
//            Request request = chain.request();
//            if (request.header("No-Authentication") == null){
//                SharedPreferences sharedPref = ???.getSharedPreferences(USER, Context.MODE_PRIVATE);
//                        request = request.newBuilder()
//                        .addHeader("Authorization", "JWT " + sharedPref.getString("auth_token", null))
//                        .build();
//            }
//            return chain.proceed(request);
//        }
//    }


    private void buildRetrofit(String url) {
//        SharedPreferences prefs = PreferenceManager.getDefaultSharedPreferences(context);
//        final String token = prefs.getString("auth_token","testtoken");
//
//        OkHttpClient client = new OkHttpClient.Builder().addInterceptor(new Interceptor() {
//            @Override
//            public Response intercept(Chain chain) throws IOException {
//                Request newRequest  = chain.request().newBuilder()
//                        .addHeader("Authorization", "Bearer " + token)
//                        .build();
//                return chain.proceed(newRequest);
//            }
//        }).build();

        Retrofit retrofit = new Retrofit.Builder()
//                .client(client)
                .baseUrl(url)
                .addConverterFactory(GsonConverterFactory.create())
                .build();

        this.github = retrofit.create(GitHub.class);
    }

    public GitHub getGithub() {
        return this.github;
    }
}