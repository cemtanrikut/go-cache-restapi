# go-cache-restapi

 Proje In memory key-value store rest api projesidir.  

 Go dilinin standart kütüphaneleri kullanılarak bir key-value memory cache uygulaması geliştirilmiştir.  

 Uygulama otuz dakikada bir cache'i okur, ve /tmp/TIMESTAMP-Data.gob dosyası içine dataları kaydeder.  
 Uygulama ayağa her kalkışında bu dosyanın içini okur ve dosya içerisindeki dataları cache'e yazar.  

 Atılan her Http Request'ler /tmp/httpServerLog.txt dosyası içerisine kaydedilir.  

 /tmp altında bu dosyaların olmaması sorun değil, eğer dosya mevcut değil ise dosyaları create eder ve yapması gereken işine devam eder.  

 Bu RestAPI aynı zamanda **HEROKU**'ya deploy edilmiştir.  
 https://go-cache-rest-api.herokuapp.com  

## Endpoints  
- ### **SET** 
 Key-Value değerini set eder. Bir BODY'e ihtiyaç vardır  

 https://go-cache-rest-api.herokuapp.com/set  
 localhost/set (local'de çalıştırmak için)  

 Örnek Body Json olacak şekilde,  
```javascript
{
    "Key" : "1",
    "Value" : {
        "Name" : "Cem",
        "Summary" : "Test the post method"
    }
}
```
    

- ### **GET** 
 Cache'deki bütün Key-Value'ları getirir.  

 https://go-cache-rest-api.herokuapp.com/get  
 localhost/get (local'de çalıştırmak için)  

- ### **GET/{Key}**  
 Key değeri string bir değerdir. Cache'de o Key'e sahip Key-Value modelini döner.  

 https://go-cache-rest-api.herokuapp.com/get/1  
 localhost/get/1 (local'de çalıştırmak için)  

- ### **DELETE/{Key}**  
 Key değeri string bir değerdir. Cache'de o Key'e sahip Key-Value modelini cacheden siler.  

 https://go-cache-rest-api.herokuapp.com/delete/1  
 localhost/delete/1 (local'de çalıştırmak için)  

- ### **FLUSH**  
 Tüm cache'i boşaltır  

 https://go-cache-rest-api.herokuapp.com/flush  
 localhost/flush (local'de çalıştırmak için)  



