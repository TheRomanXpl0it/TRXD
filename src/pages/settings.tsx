"use client"
 
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
 
import { Button } from "@/components/ui/button"
import { DatePickerWithRange } from "@/components/ui/dataRangePicker"
import { Switch } from "@/components/ui/switch"
import { Textarea } from "@/components/ui/textarea"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
} from "@/components/ui/form"
import { useContext } from "react"
import SettingsContext from "@/context/SettingsProvider"
import { toast } from "sonner"

const generateSchema = (settings: any) => {
  const schema: any = {};
  Object.keys(settings).forEach((key) => {
    settings[key].forEach((item: any) => {
      if (item.type === Boolean) {
        schema[item.title] = z.boolean().optional();
      } else if (item.type === String) {
        schema[item.title] = z.string().optional();
      } else if (item.type === Date) {
        schema[item.title] = z.date().optional();
      } else {
        schema[item.title] = z.any().optional();
      }
    });
  });
  return z.object(schema);
};


export function Settings (){
    
    const { settings, setSettings } = useContext(SettingsContext);
    const FormSchema = generateSchema(settings);
    type SettingsKeys = keyof typeof settings;

    const defaultValues = Object.keys(settings).reduce((acc: Record<string, any>, key) => {
        settings[key as SettingsKeys].forEach((item: any) => {
            acc[item.title] = item.type === Date ? new Date(item.value) : item.value;
        });
        return acc;
    }, {});

    const form = useForm<z.infer<typeof FormSchema>>({
        resolver: zodResolver(FormSchema),
        defaultValues,
    })
     
    function onSubmit(data: z.infer<typeof FormSchema>) {
        setSettings((prevSettings: any) => {
            const updatedSettings = { ...prevSettings };
            Object.keys(data).forEach((key) => {
                Object.keys(updatedSettings).forEach((settingKey) => {
                    updatedSettings[settingKey].forEach((item: any) => {
                        if (item.title === key) {
                            item.value = data[key];
                        }
                    });
                });
            });
            return updatedSettings;
        });
        toast("Settings updated");
    }

    function correctFormType(type: any, field: any){
        if(type === Boolean){
            return <Switch
                checked={field.value}
                onCheckedChange={field.onChange}
                {...field}
            />
        } else if(type === String){
            return <Textarea placeholder={field.value} {...field}/>
        } else if(type === Date){
            return <DatePickerWithRange {...field}/>
        }
        else {
            return <input
                className="w-full p-2 border rounded-lg"
                {...field}
            />
        }
    }

    function displaySettings(control: any){
        return Object.keys(settings).map((key, index) => (
            <div key={index}>
                <h1 className="scroll-m-20 text-2xl font-semibold tracking-tight mb-4 mt-4">{key}</h1>
                <div className="space-y-4">
                    {settings[key as SettingsKeys].map((item, index) => (
                        <FormField
                            key={index}
                            control={control}
                            name={item.title}
                            render={({ field }) => (
                                <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                                    <div className="space-y-0.5">
                                        <FormLabel className="text-base">{item.title}</FormLabel>
                                        <FormDescription>{item.description}</FormDescription>
                                    </div>
                                    <FormControl>
                                        {correctFormType(item.type,field)}
                                    </FormControl>
                                </FormItem>
                            )}
                        />
                    ))
                }
                </div>
            </div>
        ));
    }

    

    return (
    <>
        <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
            Settings
        </h2>
        <blockquote className="mt-6 border-l-2 pl-6 italic">
            "He who controls others may be powerful, but he who has mastered himself is mightier still."
        </blockquote>

        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="w-full space-y-6">
                {displaySettings(form.control)}
                <Button type="submit">Update Settings</Button>
            </form>
        </Form>
    </>
);
}